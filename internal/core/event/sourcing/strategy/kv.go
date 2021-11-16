package strategy

import (
	"context"
	"database/sql"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/model"
	"github.com/cacos-group/cacos/internal/core/metadata"
)

type KV struct {
	*strategy
	db *sql.DB
}

func NewKV(strategy *strategy, db *sql.DB) *KV {
	n := &KV{}
	n.strategy = strategy
	n.db = db

	return n
}

func (s *KV) GeneratorEvents(ctx context.Context, mds metadata.Metadatas) (list []model.Event) {
	key := mds.Get(metadata.Key)
	val := mds.Get(metadata.Val)

	list = []model.Event{
		model.NewKVPutEvent(key, val),
	}
	return
}
