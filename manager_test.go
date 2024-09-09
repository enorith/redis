package redis_test

import (
	"context"
	"testing"

	rds "github.com/enorith/redis"
	"github.com/redis/go-redis/v9"
)

func TestManager(t *testing.T) {
	m := rds.NewManager()
	m.Register("default", func() redis.UniversalClient {
		return redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: []string{"127.0.0.1:6379"},
		})
	})

	c, e := m.GetConnection()
	if e != nil {
		t.Error(e)
		t.Fail()
	}
	cmd := c.Ping(context.Background())
	if cmd.Val() != "PONG" {
		t.Log(cmd.String())
		t.Fail()
	}
}
