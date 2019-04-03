# oredis
oredis is a init for redigo with pool

## Redis Config
```
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

example

config := &Config{
    Network       : "tcp",
    Host          : "127.0.0.1",
    Port          : 6379,
    Password      : "",
    DB            : 0,
    Timeout       : 5,
    MaxActive     : 100,
    MaxIdle       : 50,
    MaxIdleTimeout: 5,
    Wait          : true
}
```

## GetInstance function
```
func GetInstance(config Config) redis.Conn {
}
```

## Usage
```
    redis := oredis.GetInstance(config)
```