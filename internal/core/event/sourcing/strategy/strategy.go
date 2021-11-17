package strategy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/model"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/qparams"
	"github.com/cacos-group/cacos/internal/core/metadata"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"strings"
	"time"
)

type strategy struct {
	db               *sql.DB
	etcd             *clientv3.Client
	eventSourcingMap map[model.EventSourcingName]EventSourcing
}

type Strategy interface {
	Handler(ctx context.Context, esn model.EventSourcingName, mds metadata.Metadatas) error
}

type EventSourcing interface {
	GeneratorEvents(ctx context.Context, mds metadata.Metadatas) []model.Event
}

func New(db *sql.DB, etcd *clientv3.Client) Strategy {
	s := &strategy{
		db:               db,
		etcd:             etcd,
		eventSourcingMap: make(map[model.EventSourcingName]EventSourcing),
	}
	s.eventSourcingMap[model.AddNamespace] = NewNamespace(s, db)
	s.eventSourcingMap[model.AddAppid] = NewAppid(s, db)
	s.eventSourcingMap[model.AddKV] = NewKV(s, db)

	return s
}

func (s *strategy) getStrategy(esn model.EventSourcingName) (EventSourcing, error) {
	es, ok := s.eventSourcingMap[esn]
	if !ok {
		return nil, errors.New("event sourcing undefined")
	}

	return es.(EventSourcing), nil
}

func (s *strategy) Handler(ctx context.Context, esn model.EventSourcingName, mds metadata.Metadatas) (err error) {
	es, err := s.getStrategy(esn)
	if err != nil {
		return
	}

	events := es.GeneratorEvents(ctx, mds)

	tableName := GenTableName(mds.Get(metadata.Namespace), 0)

	err = s.wal(ctx, GenRedoTableName(tableName), events)
	if err != nil {
		return
	}

	trxID := 0
	offset, err := s.Published(ctx, events)
	if err == nil {
		err = s.commit(ctx, tableName, trxID)
		if err != nil {
			return err
		}
		return
	}

	newOffset, replayedErr := s.Replayed(ctx, tableName, events, offset)
	if replayedErr == nil {
		err = s.commit(ctx, tableName, trxID)
		if err != nil {
			return err
		}
		return nil
	}

	ceeErr := s.cancelExecutedEvent(ctx, events, newOffset)
	if ceeErr != nil {
		return ceeErr
	}

	err = s.rollback(ctx, tableName, trxID)
	if err != nil {
		return err
	}

	return nil
}

func (s *strategy) GeneratorEvents(ctx context.Context, mds metadata.Metadatas) (list []model.Event) {
	return
}

func (s *strategy) wal(ctx context.Context, tableName string, events []model.Event) (err error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = s.store(tx, ctx, tableName, events)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return nil
}

func (s *strategy) store(tx *sql.Tx, ctx context.Context, tableName string, events []model.Event) (err error) {
	for _, event := range events {
		if event.EventType == model.AppidPut {
			row := tx.QueryRow(fmt.Sprintf(_existsEventLogSql, tableName), event.EventType, event.Params.Encode())
			var count int
			err := row.Scan(&count)
			if err != nil {
				return err
			}
			if count > 0 {
				return errors.New("event already exists")
			}
		}

		// todo major
		_, err = tx.ExecContext(ctx, fmt.Sprintf(_addEventLogSql, tableName), event.EventType, event.Params.Encode(), "", time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *strategy) commit(ctx context.Context, tableName string, trxID int) (err error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = func() (err error) {
		_sql := "INSERT INTO %s (`trx_id`, `event_type`, `params`, `major`, `event_time`) SELECT `trx_id`, `event_type`, `params`, `major`, `event_time` FROM %s WHERE `trx_id` = ?"
		query := fmt.Sprintf(_sql, tableName, strings.Replace(tableName, "event_log", "redo_log", 1))
		_, err = tx.ExecContext(ctx, query, trxID)
		if err != nil {
			return
		}

		_sql2 := "DELETE FROM %s WHERE `trx_id` = ?"
		query2 := fmt.Sprintf(_sql2, strings.Replace(tableName, "event_log", "redo_log", 1))
		_, err = tx.ExecContext(ctx, query2, trxID)
		if err != nil {
			return
		}

		return
	}()

	if err != nil {
		_ = tx.Rollback()
		return
	}

	_ = tx.Commit()

	return
}

func (s *strategy) Published(ctx context.Context, events []model.Event) (offset int, err error) {
	for _, event := range events {
		err = s.execEvent(ctx, event)
		if err != nil {
			break
		}
		offset += 1
	}

	if err != nil {
		return
	}

	return
}

func (s *strategy) execEvent(ctx context.Context, event model.Event) (err error) {
	switch event.EventType {
	case model.InfoNamespacePut, model.InfoAppidPut:
		_, err = s.etcd.KV.Put(ctx, event.Params.Get(qparams.Key), event.Params.Get(qparams.Val))
	case model.AppidPut, model.KVPut:
		_, err = s.etcd.KV.Put(ctx, event.Params.Get(qparams.Key), event.Params.Get(qparams.Val))
	case model.UserAdd:
		_, err = s.etcd.UserAdd(ctx, event.Params.Get(qparams.User), event.Params.Get(qparams.Password))
	case model.RoleAdd:
		_, err = s.etcd.RoleAdd(ctx, event.Params.Get(qparams.Role))
	case model.UserGrantRole:
		_, err = s.etcd.UserGrantRole(ctx, event.Params.Get(qparams.User), event.Params.Get(qparams.Role))
	case model.RoleGrantPermission:
		perm, _ := strconv.Atoi(event.Params.Get(qparams.Perm))
		_, err = s.etcd.RoleGrantPermission(ctx, event.Params.Get(qparams.Role), event.Params.Get(qparams.Key), "\\0", clientv3.PermissionType(perm))

	case model.InfoAppidDel:
		_, err = s.etcd.KV.Delete(ctx, event.Params.Get(qparams.Key))
	case model.AppidDel:
		_, err = s.etcd.KV.Delete(ctx, event.Params.Get(qparams.Key))
	case model.KVDel:
	case model.UserDel:
		_, err = s.etcd.UserDelete(ctx, event.Params.Get(qparams.User))
	case model.RoleDel:
		_, err = s.etcd.RoleDelete(ctx, event.Params.Get(qparams.Role))

	default:
		err = errors.New("event_type undefined")
	}

	return
}

func (s *strategy) Replayed(ctx context.Context, tableName string, events []model.Event, offset int) (newOffset int, err error) {
	//todo retry次数
	for offset < len(events) {
		err = s.execEvent(ctx, events[offset])
		if err != nil {
			break
		}

		offset++
	}

	newOffset = offset

	if err != nil {
		return
	}

	return
}

func (s *strategy) cancelExecutedEvent(ctx context.Context, events []model.Event, offset int) (err error) {
	// 执行补偿事件
	successEvents := events[:offset]
	for _, event := range successEvents {
		if event.Cancel.EventType != "" {
			err = s.execEvent(ctx, event.Cancel.Format())
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		//todo 计划任务恢复
		return
	}
	return
}

func (s *strategy) rollback(ctx context.Context, tableName string, trxID int) (err error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}

	_sql2 := "DELETE FROM %s WHERE `trx_id` = ?"
	query2 := fmt.Sprintf(_sql2, strings.Replace(tableName, "event_log", "redo_log", 1))
	_, err = conn.ExecContext(ctx, query2, trxID)
	if err != nil {
		return
	}

	return
}
