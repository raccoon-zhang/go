package dbPool

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type Pool struct {
	busyPool    sync.Map
	freePool    sync.Map
	isDestroyed bool
	masterOps   *redis.Options
}

func InitPool(driver, dsn string, ops *redis.Options, size int) (*Pool, error) {
	var pool = &Pool{}
	pool.masterOps = ops
	for count := 0; count < size; count++ {
		db, err := sql.Open(driver, dsn)
		if err != nil {
			fmt.Println(err)
			//清除所有的已经插入的数据
			pool.freePool = sync.Map{}
			return nil, errors.New("create error")
		}
		pool.freePool.Store(db, db)
	}
	return pool, nil
}

func (pool Pool) NewDb() (db *sql.DB, err error) {
	pool.freePool.Range(func(key, value any) bool {
		if db == nil {
			db, _ = value.(*sql.DB)
			return false
		}
		return true
	})
	if db != nil {
		pool.freePool.Delete(db)
		pool.busyPool.Store(db, db)
		return db, nil
	} else {
		return nil, errors.New("no db is free")
	}
}
func (pool Pool) DeleteDb(db *sql.DB) {
	if pool.isDestroyed {
		db.Close()
		return
	}
	pool.busyPool.Delete(db)
	pool.freePool.Store(db, db)
}

func (pool Pool) DestroyPool() {
	pool.busyPool.Range(func(key, value any) bool {
		db, ok := value.(*sql.DB)
		if ok {
			db.Close()
		}
		pool.busyPool.Delete(key)
		return true
	})
	pool.freePool.Range(func(key, value any) bool {
		db, ok := value.(*sql.DB)
		if ok {
			db.Close()
		}
		pool.freePool.Delete(key)
		return true
	})
	pool.isDestroyed = true
}

func (pool Pool) NewRedisCliForWrite() *redis.Client {
	return redis.NewClient(pool.masterOps)
}

func (pool Pool) NewRedisCliForRead(ops *redis.Options) *redis.Client {
	return redis.NewClient(ops)
}

func (pool Pool) DeleteRedisCli(cli *redis.Client) {
	cli.Close()
}
