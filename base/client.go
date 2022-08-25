package base

import (
	"entgo.io/ent/examples/fs/ent"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	DB *ent.Client
	RE *redis.Client
}
