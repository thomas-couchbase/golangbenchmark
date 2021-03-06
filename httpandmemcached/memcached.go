package main

import (
	"fmt"
	"github.com/kklis/gomemcache"
	//"sync/atomic"
)

type Pool struct {
	stack chan *gomemcache.Memcache
	addr  string // the memcached server address
	count int32 // The total count of the memcached connections
}

func New(addr string, size int) *Pool {
	pool := new(Pool)
	pool.addr = addr
	pool.stack = make(chan *gomemcache.Memcache, size)
	for i := 0; i < size; i++ {
		c, err := gomemcache.Dial(addr)
		if err != nil {
			fmt.Printf("Connect to %v failed : %v\n", addr, err.Error())
			return nil
		}
		pool.stack <- c
		//atomic.AddInt32(&pool.count, 1)
	}
	return pool
}

func (pool *Pool) Get() *gomemcache.Memcache {
	select {
	case conn := <-pool.stack:
		//atomic.AddInt32(&pool.count, 1)
		//fmt.Printf("get memcached conn from pool: %v, pool size %v\n", conn, len(pool.stack));
		return conn
	default:
		conn, err := gomemcache.Dial(pool.addr)
		if err != nil {
			fmt.Printf("Connect to %v failed : %v\n", pool.addr, err.Error())
			return nil
		}
		//atomic.AddInt32(&pool.count, 1)
		//fmt.Printf("create a new conn %v, pool size %v\n", conn, len(pool.stack));
		return conn
	}
}

func (pool *Pool) Put(conn *gomemcache.Memcache) {
	select {
	case pool.stack <- conn:
		//atomic.AddInt32(&pool.count, -1)
		//fmt.Printf("put conn %v to pool, pool size %v\n", conn, len(pool.stack));
	default:
	}
}
