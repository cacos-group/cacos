package strategy

import (
	"context"
	"database/sql"
	"strings"
)

type Appid struct {
	*Strategy
	db *sql.DB
}

func NewAppid(strategy *Strategy, db *sql.DB) *Appid {
	n := &Appid{}
	n.Strategy = strategy
	n.db = db

	return n
}

func (s *Appid) Prepare(ctx context.Context, namespace string, appid string) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// add appid
	_, err = conn.ExecContext(ctx, _addInfo, namespace, appid, 2)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") && strings.Contains(err.Error(), "for key 'uniq_namespace_appid'") {
			return nil
		}
		return err
	}

	return nil
}

func (s *Appid) Replayed(ctx context.Context) error {
	return nil
}
