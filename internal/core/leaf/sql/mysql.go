package sql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var RowsNotAffectedErr = fmt.Errorf("rows not affected")

// Config sql config.
type Config struct {
	DSN    string // write data source name.
	BizTag string
}

type MySQL struct {
	db   *sql.DB
	conf *Config
	step int
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(db *sql.DB, bizTag string, step int) (mysql *MySQL, cf func(), err error) {
	return &MySQL{
		db: db,
		conf: &Config{
			BizTag: bizTag,
		},
		step: step,
	}, func() {}, nil
}

func (m *MySQL) InitBizTag(ctx context.Context, bizTag string, maxID uint64, step int, description string) error {
	db := m.db

	_, err := db.Exec(_insertSql, bizTag, maxID, step, description)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) GetEndID(ctx context.Context) (startID uint64, endID uint64, step int, err error) {
	db := m.db
	tx, err := db.Begin()
	if err != nil {
		return
	}

	bizTag := m.conf.BizTag
	res, err := tx.Exec(_updateSql, bizTag)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = RowsNotAffectedErr
		return
	}

	row := tx.QueryRow(_selectSql, bizTag)

	var maxID uint64
	err = row.Scan(&maxID, &step)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	m.step = step
	startID = maxID - uint64(step) + 1
	endID = maxID
	return
}
