package data

import (
	"customer/internal/conf"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo,NewCustomerData)

// Data .
type Data struct {
	// TODO wrapped database client
	Rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{}
	//连接redis，使用服务的配置，c就是解析之后的配置信息
	redisURL := fmt.Sprintf("redis://%s/1", c.Redis.Addr) //1是数据库
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		data.Rdb = nil
	}
	data.Rdb = redis.NewClient(options) //建立客户端，不会立即连接需，要执行命令时才会连接

	cleanup := func() {
		//清理了 Redis 连接
		_ = data.Rdb.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return data, cleanup, nil
}
