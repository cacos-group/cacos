package leaf

import (
	"github.com/cacos-group/cacos/internal/core/leaf/sql"
	"context"
	xsql "database/sql"
	"fmt"
	"log"
	"sync"
)

type idPool struct {
	NextID uint64
	MaxID  uint64
	Step   int
}

type Leaf interface {
	NextId() (uint64, error)
}

type leaf struct {
	mu           sync.Mutex
	idPool1Mutex sync.Mutex

	idPool0 *idPool
	idPool1 *idPool
	ch      chan bool

	mysql *sql.MySQL

	config *Config
}

type Config struct {
	DSN    string // write data source name.
	BizTag string
}

func New(db *xsql.DB) (l Leaf, cf func(), err error) {
	return newLeaf(db)
}

func newLeaf(db *xsql.DB) (l *leaf, cf func(), err error) {
	bizType := "cacos"
	step := 100

	ch := make(chan bool, step)

	ms, _, err := sql.NewMySQL(db, bizType, step)
	if err != nil {
		return
	}

	err = ms.InitBizTag(context.TODO(), bizType, 1, step, "cacos")
	if err != nil {
		return
	}

	l = &leaf{
		ch:    ch,
		mysql: ms,
	}

	err = l.initIdPool0()
	if err != nil {
		return
	}

	l.initDaemon()

	cf = func() {

	}

	return
}

func (l *leaf) NextId() (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.getNextId()
}

var IDPool0EmptyErr = fmt.Errorf("idPool0 empty")

func (l *leaf) getNextId() (uint64, error) {
	idPool := l.idPool0
	if idPool == nil {
		return 0, IDPool0EmptyErr
	}

	nextID := idPool.NextID
	if nextID <= idPool.MaxID {
		l.idPool0.NextID = nextID + 1
		l.ch <- true
		return nextID, nil
	}

	if l.idPool1 == nil {
		err := l.initIdPool1()
		if err != nil {
			return 0, err
		}
	}

	l.idPool0 = l.idPool1
	l.idPool1 = nil

	return l.getNextId()
}

func (l *leaf) initDaemon() {
	go func() {
		for {
			select {
			case <-l.ch:
				err := l.initIdPool1()
				if err != nil {
					log.Printf("daemon initIdPool1 err: %v", err)
				}
			}
		}
	}()
}

func (l *leaf) initIdPool0() error {
	startID, endID, step, err := l.mysql.GetEndID(context.TODO())
	if err != nil {
		return err
	}

	maxID := endID
	currentID := startID

	l.idPool0 = &idPool{
		NextID: currentID,
		MaxID:  maxID,
		Step:   step,
	}

	return nil
}

func (p *idPool) left() uint64 {
	if p.NextID > p.MaxID {
		return 0
	}
	return p.MaxID - p.NextID + 1
}

func (p *idPool) total() uint64 {
	return uint64(p.Step)
}

func (l *leaf) initIdPool1() error {
	l.idPool1Mutex.Lock()
	defer l.idPool1Mutex.Unlock()

	if l.idPool1 != nil && l.idPool1.MaxID > 0 {
		return nil
	}

	if l.idPool0 != nil {
		rate := float64(l.idPool0.left()) / float64(l.idPool0.total())
		if rate > 0.6 {
			return nil
		}
	}

	startID, endID, step, err := l.mysql.GetEndID(context.TODO())
	if err != nil {
		return err
	}

	maxID := endID
	currentID := startID

	l.idPool1 = &idPool{
		NextID: currentID,
		MaxID:  maxID,
		Step:   step,
	}

	return nil
}
