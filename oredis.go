// Package oredis is own mysql
package oredis

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
)

// Config Reids 配置
type Config struct {
	Network        string // TCP
	Host           string
	Port           int
	Password       string
	DB             int
	Timeout        int
	MaxActive      int  // 最大连接数，即最多的tcp连接数
	MaxIdle        int  // 最大空闲连接数，即会有这么多个连接提前等待着
	MaxIdleTimeout int  // 空闲连接超时时间
	Wait           bool // 如果超过最大连接，是报错，还是等待。true 为等待，false为报错
}

// 包内变量，存储实例相关对象
var packageOnce = map[string]*sync.Once{}
var packageInstance = map[string]*redis.Pool{}
var packageMutex = &sync.Mutex{}

// Pool Redis 连接池
func pool(config Config) *redis.Pool {
	// fmt.Println("pool init")
	// 建立连接池
	pool := &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: time.Duration(config.MaxIdleTimeout) * time.Second,
		Wait:        config.Wait,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial(config.Network, fmt.Sprintf("%s:%d", config.Host, config.Port),
				redis.DialPassword(config.Password),
				redis.DialDatabase(config.DB),
				redis.DialConnectTimeout(time.Duration(config.Timeout)*time.Second),
				redis.DialReadTimeout(time.Duration(config.Timeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(config.Timeout)*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}

	return pool
}

// GetPoolInstance 根据配置信息初始化 只初始化一次
func GetPoolInstance(config Config) *redis.Pool {
	packageMutex.Lock()
	defer packageMutex.Unlock()

	md5byte := md5.Sum([]byte(fmt.Sprintf("%s%s%d%s%d", config.Network, config.Host, config.Port, config.Password, config.DB)))
	md5key := fmt.Sprintf("%x", md5byte)

	// 如果有值直接返回
	if v, ok := packageInstance[md5key]; ok {
		// fmt.Println("direct")
		return v
	}

	// 如果once 不存在
	if _, ok := packageOnce[md5key]; !ok {
		var once = &sync.Once{}
		var obj *redis.Pool
		// var err error
		once.Do(func() {
			obj = pool(config)

			packageInstance[md5key] = obj
			packageOnce[md5key] = once
			// fmt.Printf("init %p %v\n", obj, obj)
		})

		return obj
	}

	return nil
}

// GetInstance get a redis instance from redis pool
func GetInstance(config Config) redis.Conn {
	// 初始化连接池
	redisPool := GetPoolInstance(config)

	// 从池里获取连接
	rc := redisPool.Get()
	// defer rc.Close()

	if nil == rc.Err() {
		return rc
	}

	// 失败重连 connection reset by peer
	retryTimes := 10
	for i := 0; i < retryTimes; i++ {
		time.Sleep(50 * time.Millisecond)

		rc = GetInstance(config)
		if nil == rc.Err() {
			return rc
		}
	}

	// 错误的 rc 也返回
	return rc
}

// GetInstancePanic panic when error occurred
func GetInstancePanic(config Config) redis.Conn {
	// 初始化连接池
	redisPool := GetPoolInstance(config)

	// 从池里获取连接
	rc := redisPool.Get()
	// defer rc.Close()

	if nil != rc.Err() {
		panic("Redis " + rc.Err().Error())
	}

	return rc
}
