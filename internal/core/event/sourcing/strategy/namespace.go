package strategy

import (
	"context"
	"database/sql"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/model"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/qparams"
	"github.com/cacos-group/cacos/internal/core/metadata"
)

type Namespace struct {
	*strategy
	db *sql.DB
}

func NewNamespace(strategy *strategy, db *sql.DB) *Namespace {
	n := &Namespace{}
	n.strategy = strategy
	n.db = db

	return n
}

func (s *Namespace) GeneratorEvents(ctx context.Context, mds metadata.Metadatas) (list []model.Event) {
	params := qparams.Metadatas{}
	params.Set(qparams.Key, model.GenNamespaceKey(mds.Get(metadata.Namespace)))
	params.Set(qparams.Val, mds.Get(metadata.Namespace))

	list = []model.Event{
		{
			EventType: model.InfoNamespacePut,
			Params:    params,
		},
	}
	return
}
