package dbPool

import (
	"database/sql"
	"errors"
	"log/slog"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type Pool struct {
	busyPool    sync.Map
	freePool    sync.Map
	isDestroyed bool
	failoverOps *redis.FailoverOptions
}

func InitPool(driver, dsn string, ops *redis.FailoverOptions, size int) (*Pool, error) {
	var pool = &Pool{}
	pool.failoverOps = ops
	for count := 0; count < size; count++ {
		db, err := sql.Open(driver, dsn)
		if err != nil {
			slog.Error(err.Error())
			//清除所有的已经插入的数据
			pool.freePool = sync.Map{}
			return nil, errors.New("create error")
		}
		pool.freePool.Store(db, db)
	}
	return pool, nil
}

func (pool *Pool) NewDb() (db *sql.DB, err error) {
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
func (pool *Pool) DeleteDb(db *sql.DB) {
	if pool.isDestroyed {
		db.Close()
		return
	}
	pool.busyPool.Delete(db)
	pool.freePool.Store(db, db)
}

func (pool *Pool) DestroyPool() {
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

func (pool *Pool) NewRedisCliForWrite() *redis.Client {
	//获取主节点用于写操作
	client := redis.NewFailoverClient(pool.failoverOps)
	return client
}

func (pool *Pool) NewRedisCliForRead() *redis.ClusterClient {
	// 创建Redis哨兵客户端实例
	client := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:    pool.failoverOps.MasterName,    // 主节点的名称
		SentinelAddrs: pool.failoverOps.SentinelAddrs, // 哨兵节点的地址列表
		RouteRandomly: true,                           //随机节点
	})
	return client
}

func (pool *Pool) DeleteRedisCli(cli any) {
	if val, ok := cli.(*redis.Client); ok {
		val.Close()
	} else if val, ok := cli.(*redis.ClusterClient); ok {
		val.Close()
	} else {
		slog.Warn("you passed a wrong client")
	}
}
