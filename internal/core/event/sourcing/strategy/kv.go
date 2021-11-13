package strategy

import (
	"context"
	"database/sql"
)

type KV struct {
	*Strategy
	db *sql.DB
}

func NewKV(strategy *Strategy, db *sql.DB) *KV {
	n := &KV{}
	n.Strategy = strategy
	n.db = db

	return n
}

func (s *KV) Prepare(ctx context.Context, namespace string, appid string) error {
	return nil
}

func (s *KV) Replayed(ctx context.Context) error {
	return nil
}
