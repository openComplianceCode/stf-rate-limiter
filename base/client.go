package base

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	DB *sql.DB
	RE *redis.Client
}
