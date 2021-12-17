package controller

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/model"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/pb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type controller struct {
	db   *sql.DB
	etcd *clientv3.Client
}

type Controller interface {
	Handler(ctx context.Context, req *model.HandlerReq) error
}

func New(db *sql.DB, etcd *clientv3.Client) Controller {
	c := &controller{
		db:   db,
		etcd: etcd,
	}

	return c
}

func (s *controller) Handler(ctx context.Context, req *model.HandlerReq) (err error) {
	tableName := GenTableName(req.Namespace, 0)

	err = s.wal(ctx, GenWalLogTableName(tableName), req.Events)
	if err != nil {
		return
	}

	trxID := 0
	offset, err := s.Published(ctx, req.Events)
	if err == nil {
		err = s.commit(ctx, tableName, trxID)
		if err != nil {
			return err
		}
		return
	}

	newOffset, replayedErr := s.Replayed(ctx, tableName, req.Events, offset)
	if replayedErr == nil {
		err = s.commit(ctx, tableName, trxID)
		if err != nil {
			return err
		}
		return nil
	}

	ceeErr := s.cancelExecutedEvent(ctx, req.Events, newOffset)
	if ceeErr != nil {
		return ceeErr
	}

	rErr := s.rollback(ctx, tableName, trxID)
	if rErr != nil {
		return rErr
	}

	return
}

func (s *controller) wal(ctx context.Context, tableName string, events []*pb.Event) (err error) {
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
		fmt.Println(err)
		tx.Rollback()
		return
	}

	tx.Commit()

	return nil
}

func (s *controller) store(tx *sql.Tx, ctx context.Context, tableName string, events []*pb.Event) (err error) {
	for _, event := range events {
		switch event.Payload.(type) {
		case *pb.Event_ServicePut:
			row := tx.QueryRow(fmt.Sprintf(_existsEventLogSql, tableName), model.MustMarshal(event))
			var count int
			err := row.Scan(&count)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if count > 0 {
				return errors.New("event already exists")
			}
		default:
			// todo major
			_, err = tx.ExecContext(ctx, fmt.Sprintf(_addWalLogSql, tableName), model.MustMarshal(event), "", time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}

func (s *controller) commit(ctx context.Context, tableName string, trxID int) (err error) {
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
		_sql := _addEventLogSql
		query := fmt.Sprintf(_sql, tableName, GenWalLogTableName(tableName))
		_, err = tx.ExecContext(ctx, query, trxID)
		if err != nil {
			fmt.Println(err)
			return
		}

		_sql2 := _delWalLogSql
		query2 := fmt.Sprintf(_sql2, GenWalLogTableName(tableName))
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

func (s *controller) Published(ctx context.Context, events []*pb.Event) (offset int, err error) {
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

func (s *controller) execEvent(ctx context.Context, event *pb.Event) (err error) {
	switch event.Payload.(type) {
	case *pb.Event_NamespacePut:
		_, err = s.etcd.KV.Put(ctx, event.GetNamespacePut().Key, event.GetNamespacePut().Val)
	case *pb.Event_ServicePut:
		_, err = s.etcd.KV.Put(ctx, event.GetServicePut().Key, event.GetServicePut().Val)
	case *pb.Event_KvPut:
		_, err = s.etcd.KV.Put(ctx, event.GetKvPut().Key, event.GetKvPut().Val)
	case *pb.Event_UserAdd:
		_, err = s.etcd.UserAdd(ctx, event.GetUserAdd().User, event.GetUserAdd().Password)
	case *pb.Event_RoleAdd:
		_, err = s.etcd.RoleAdd(ctx, event.GetRoleAdd().Role)
	case *pb.Event_UserGrantRole:
		_, err = s.etcd.UserGrantRole(ctx, event.GetUserGrantRole().User, event.GetUserGrantRole().Role)
	case *pb.Event_RoleGrantPermission:
		_, err = s.etcd.RoleGrantPermission(ctx,
			event.GetRoleGrantPermission().Role,
			event.GetRoleGrantPermission().Key,
			"\\0",
			clientv3.PermissionType(event.GetRoleGrantPermission().Perm))

	case *pb.Event_NamespaceDel:
		_, err = s.etcd.KV.Delete(ctx, event.GetNamespaceDel().Key)
	case *pb.Event_ServiceDel:
		_, err = s.etcd.KV.Delete(ctx, event.GetServiceDel().Key)
	//case model.KVDel:
	case *pb.Event_UserDel:
		_, err = s.etcd.UserDelete(ctx, event.GetUserDel().User)
	case *pb.Event_RoleDel:
		_, err = s.etcd.RoleDelete(ctx, event.GetRoleDel().Role)

	case nil:

	default:
		err = errors.New("event_type undefined")
	}

	return
}

func (s *controller) Replayed(ctx context.Context, tableName string, events []*pb.Event, offset int) (newOffset int, err error) {
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

func (s *controller) cancelExecutedEvent(ctx context.Context, events []*pb.Event, offset int) (err error) {
	// 执行补偿事件
	//successEvents := Events[:offset]
	//for _, event := range successEvents {
	//	if event.Cancel.EventType != "" {
	//		err = s.execEvent(ctx, event.Cancel.Format())
	//		if err != nil {
	//			break
	//		}
	//	}
	//}
	//if err != nil {
	//	//todo 计划任务恢复
	//	return
	//}
	return
}

func (s *controller) rollback(ctx context.Context, tableName string, trxID int) (err error) {
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
