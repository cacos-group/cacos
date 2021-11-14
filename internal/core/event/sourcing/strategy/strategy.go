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
)

type Strategy struct {
	db   *sql.DB
	etcd *clientv3.Client
}

func NewStrategy(db *sql.DB, etcd *clientv3.Client) *Strategy {
	return &Strategy{
		db:   db,
		etcd: etcd,
	}
}

func (s *Strategy) GeneratorEvents(ctx context.Context, mds metadata.Metadatas) (list []model.Event) {
	return
}

func (s *Strategy) Presentation(ctx context.Context, tableName string, events []model.Event) (err error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = func() error {
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

			_, err = tx.ExecContext(ctx, fmt.Sprintf(_addEventLogSql, tableName), event.EventType, event.Params.Encode())
			if err != nil {
				return err
			}
		}
		return nil
	}()

	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return nil
}

func (s *Strategy) Published(ctx context.Context, events []model.Event) error {
	for _, event := range events {
		var err error
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
		default:
			err = errors.New("event_type undefined")
		}
		if err != nil {
			return err
		}
	}

	return nil
}
