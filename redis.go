package goapis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/lalolv/goutil"
)

// RedisDo 单条 redis 服务器执行
// @cmd 命令名称
// @args 参数
func RedisDo(redisPool *redis.Pool, dbIdx int, cmd string, args ...interface{}) interface{} {
	// 参数
	opt := []interface{}{}
	for _, arg := range args {
		opt = append(opt, arg)
	}
	// redis 连接
	cli := redisPool.Get()
	//选定数据库索引
	cli.Do("SELECT", dbIdx)
	// 执行命令
	r, err := cli.Do(cmd, opt...)
	if err != nil {
		// 错误处理
		fmt.Println(err.Error())
	}
	return r
}

// RedisDos 批量 redis 服务器执行
// 二维数组。每条数组一个命令。
// 形如： [["c1", "f1", "v1"], ["c2", "f2", "v2", "t2"]]
// @cmds 批量命令
func RedisDos(redisPool *redis.Pool, dbIdx int, cmds [][]interface{}) []interface{} {
	// redis 连接
	cli := redisPool.Get()
	//选择数据库
	cli.Do("SELECT", dbIdx)
	//执行命令
	var results []interface{}
	for _, ele := range cmds {
		r, err := cli.Do(ele[0].(string), ele[1:]...)
		results = append(results, r)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return results
}

// RedisDoInt 执行redis命令返回int类型的数据，如果遇到错误则返回0
// @dbIndex
// @cmd
// @args
func RedisDoInt(redisPool *redis.Pool, dbIdx int, cmd string, args ...interface{}) int {
	ret := RedisDo(redisPool, dbIdx, cmd, args...)
	if t, ok := ret.([]uint8); ok {
		return goutil.U8sInt(t)
	} else if t, ok := ret.(int64); ok {
		return int(t)
	}
	return 0
}

// NewPool 建立redisPool连接池
// @server
// @password 如果为空就不输入密码
func NewPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     20,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err = c.Do("AUTH", password); err != nil {
					c.Close()
					fmt.Println("Redis AUTH:", err.Error())
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			fmt.Println("Redis Ping:", err.Error())
			return err
		},
	}
}
