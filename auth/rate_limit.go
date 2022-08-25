package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

var (
	GeneralRL = map[string]int{
		"nologin":  10000,
		"everyone": 40000,
		"gold":     160000,
		"diamon":   640000,
		"root":     10000000,
	}
)

func CheckRateLimit(ctx context.Context, reClient *redis.Client, n int) bool {
	remoteIP := ForIpContext(ctx)
	userInfo := ForContext(ctx)
	limiter := redis_rate.NewLimiter(reClient)
	if userInfo == nil { // no login check IP
		res, err := limiter.AllowN(ctx, "IP_RL:"+remoteIP, redis_rate.Limit{
			Rate:   GeneralRL["nologin"],
			Period: time.Hour,
			Burst:  GeneralRL["nologin"] * 24,
		}, n)
		if err != nil {
			log.Println("limiter checking err: ", err)
			return false
		}
		return res.Allowed >= n
	} else {
		res, err := limiter.AllowN(ctx, "USER_RL:"+fmt.Sprint(userInfo.UserId), redis_rate.Limit{
			Rate:   GeneralRL[userInfo.Role],
			Period: time.Hour,
			Burst:  GeneralRL[userInfo.Role] * 24,
		}, n)

		if err != nil {
			log.Println("limiter checking err: ", err)
			return false
		}
		return res.Allowed >= n
	}
}
