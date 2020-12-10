package sharding

import (
	"database/sql"
)

const (
	shardCntMax = 100
)

type DBShardPool struct {
	shardList []*sql.DB
}

func NewDBShardPool() *DBShardPool {
	return &DBShardPool{}
}

func (dbp *DBShardPool) AddShard(s *sql.DB) {
	dbp.shardList = append(dbp.shardList, s)
}

func (dbp *DBShardPool) GetShard(dialogID int) *sql.DB {
	shardCnt := len(dbp.shardList)
	if shardCnt == 1 {
		return dbp.shardList[0]
	}

	d := dialogID % shardCntMax
	var k int = (shardCntMax / shardCnt)

	for i := 0; i < shardCnt; i++ {
		if d <= (i+1)*k {
			return dbp.shardList[i]
		}
	}

	return nil
}

func (dbp *DBShardPool) Close() error {
	var err error
	for _, sConn := range dbp.shardList {
		err = sConn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
