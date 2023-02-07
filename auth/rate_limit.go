package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alec-z/rate_limiter/base"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

func CheckRateLimit(ctx context.Context, reClient *redis.Client, n int) (bool, int) {
	remoteIP := ForIpContext(ctx)
	userInfo := ForContext(ctx)
	limiter := redis_rate.NewLimiter(reClient)
	if userInfo == nil { // no login check IP
		res, err := limiter.AllowN(ctx, "IP_RL:"+remoteIP, redis_rate.Limit{
			Rate:   base.GeneralRL["nologin"].(int),
			Period: time.Hour * 24,
			Burst:  base.GeneralRL["nologin"].(int) * 24,
		}, n)
		if err != nil {
			log.Println("limiter checking err: ", err)
			return false, res.Remaining
		}
		log.Println("remoteIP : " + remoteIP + " consume :" + fmt.Sprint(n) + " Allow :" + fmt.Sprint(res.Allowed))
		return res.Allowed >= n, res.Remaining
	} else {
		res, err := limiter.AllowN(ctx, "USER_RL:"+fmt.Sprint(userInfo.UserId), redis_rate.Limit{
			Rate:   base.GeneralRL[userInfo.Role].(int),
			Period: time.Hour * 24,
			Burst:  base.GeneralRL[userInfo.Role].(int) * 24,
		}, n)

		if err != nil {
			log.Println("limiter checking err: ", err)
			return false, res.Remaining
		}
		log.Println("userID: " + fmt.Sprint(userInfo.UserId) + " consume :" + fmt.Sprint(n) + " Allow :" + fmt.Sprint(res.Allowed))
		return res.Allowed >= n, res.Remaining
	}
}
