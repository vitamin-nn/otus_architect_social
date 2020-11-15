package db

import (
	"database/sql"
	"sync/atomic"
)

type DBReplPool struct {
	master  *sql.DB
	slaves  []*sql.DB
	rrCount uint64 // round robin counter
}

func NewDBPool(m *sql.DB) *DBReplPool {
	return &DBReplPool{
		master: m,
	}
}

func (dbp *DBReplPool) AddSlave(s *sql.DB) {
	dbp.slaves = append(dbp.slaves, s)
}

func (dbp *DBReplPool) GetMaster() *sql.DB {
	return dbp.master
}

func (dbp *DBReplPool) GetSlave() *sql.DB {
	n := len(dbp.slaves)
	if n < 1 {
		return dbp.GetMaster()
	}

	if n == 1 {
		return dbp.slaves[0]
	}

	m := int(atomic.AddUint64(&dbp.rrCount, 1) % uint64(n))
	return dbp.slaves[m]
}

func (dbp *DBReplPool) Close() error {
	err := dbp.master.Close()
	if err != nil {
		return err
	}

	for _, slConn := range dbp.slaves {
		err = slConn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
