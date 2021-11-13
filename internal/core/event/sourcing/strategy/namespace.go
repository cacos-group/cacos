package strategy

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/model"
)

type Namespace struct {
	*Strategy
	db *sql.DB
}

func NewNamespace(strategy *Strategy, db *sql.DB) *Namespace {
	n := &Namespace{}
	n.Strategy = strategy
	n.db = db

	return n
}

func (s *Namespace) Prepare(ctx context.Context, namespace string, appid string) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = func(tx *sql.Tx) (err error) {
		// add namespace
		_, err = conn.ExecContext(ctx, _addInfo, namespace, appid, 1)
		if err != nil {
			return err
		}

		tableName := GenTableName(namespace, appid)

		// create event_log_%s_%s
		_, err = tx.ExecContext(ctx, fmt.Sprintf(_createEventLogSql, tableName))
		if err != nil {
			return
		}
		return
	}(tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (s *Namespace) Presentation(ctx context.Context, tableName string, events []model.Event) error {
	return nil
}

func (s *Namespace) Published(ctx context.Context, events []model.Event) error {
	return nil
}

func (s *Namespace) Replayed(ctx context.Context) error {
	return nil
}
