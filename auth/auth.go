package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alec-z/rate_limiter/base"
	"github.com/alec-z/rate_limiter/model"
	"github.com/golang-jwt/jwt/v4"
)

type contextKey struct {
	name string
}

// A stand-in for our database backed user object
type UserInfo struct {
	UserId int64  `json:"user_id,omitempty"`
	Role   string `json:"role,omitempty"`
}

type CustomClaims struct {
	jwt.StandardClaims
	UserInfo
}

const (
	JWT_VALID_MINS = 60
)

var userCtxKey = &contextKey{"user"}
var ipCtxKey = &contextKey{"remote_ip"}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(c *base.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		saveIPContext(w, r)
		var rr *http.Request = r

		if tokenString := strings.TrimSpace(r.Header.Get("Authorization")); tokenString != "" {
			rr = checkAuthorizationForLogin(tokenString, c, w, r)
		} else if cookie, err := r.Cookie("auth-cookie"); err == nil && c != nil && strings.TrimSpace(cookie.Value) != "" {
			rr = checkCookieForLogin(cookie, w, r)
		}
		if rr != nil {
			next.ServeHTTP(w, rr)
			refreshCookie(w, rr)
		}

	})
}

func refreshCookie(w http.ResponseWriter, r *http.Request) {
	if userInfo := ForContext(r.Context()); userInfo != nil {
		var newClaim CustomClaims
		newClaim.UserInfo = *userInfo
		if newTokenString, err := CreateOrRefreshToken(&newClaim); err == nil {
			var cookie http.Cookie
			cookie.Name = "auth-cookie"
			cookie.Value = newTokenString
			cookie.MaxAge = 60 * (JWT_VALID_MINS - 5)
			http.SetCookie(w, &cookie)
		} else {
			log.Println("Refresh token failed: ", err)
		}
	}

}

func CreateOrRefreshToken(c *CustomClaims) (tokenString string, err error) {
	c.ExpiresAt = time.Now().Add(time.Minute * JWT_VALID_MINS).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString(base.HmacSecret)
	return
}

func checkAuthorizationForLogin(apiTokenStr string, c *base.Client, w http.ResponseWriter, r *http.Request) *http.Request {
	tokenStrArr := strings.Split(apiTokenStr, " ")

	if len(tokenStrArr) > 1 {
		apiTokenStr = tokenStrArr[len(tokenStrArr)-1]
	} else {
		apiTokenStr = tokenStrArr[0]
	}

	userInfo := getUserInfoFromRedis(r.Context(), c, apiTokenStr)
	if userInfo != nil && userInfo.UserId != 0 {
		return saveUserContext(userInfo, r)
	}

	http.Error(w, "Invalid Authorization Token", http.StatusUnauthorized)
	return nil
}

func checkCookieForLogin(cookie *http.Cookie, w http.ResponseWriter, r *http.Request) *http.Request {
	tokenString := strings.TrimSpace(cookie.Value)
	if claims, err := extractClaims(tokenString); err == nil && claims.UserId != 0 { // check jwt token
		rr := saveUserContext(&claims.UserInfo, r)
		if newTokenString, err := CreateOrRefreshToken(claims); err == nil {
			cookie.Value = newTokenString
			cookie.MaxAge = 60 * (JWT_VALID_MINS - 5)
			http.SetCookie(w, cookie) // refresh jwt token and set new cookie
		} else {
			log.Println("Refresh token failed: ", err)
		}
		return rr
	}
	http.Error(w, "Invalid Cookie", http.StatusUnauthorized)
	return nil
}

func saveUserContext(userInfo *UserInfo, r *http.Request) *http.Request {
	if userInfo != nil {
		ctx := context.WithValue(r.Context(), userCtxKey, userInfo)
		return r.WithContext(ctx)
	}
	return r
}

func saveIPContext(w http.ResponseWriter, r *http.Request) *http.Request {
	remoteIP := r.Header.Get(base.IP_HEADER_KEY)
	w.Header().Set("z-client-ip", remoteIP)
	ctx := context.WithValue(r.Context(), ipCtxKey, remoteIP)
	return r.WithContext(ctx)
}

func extractClaims(tokenString string) (c *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return base.HmacSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		fmt.Println(claims.UserInfo)
		return claims, nil
	} else {
		return nil, fmt.Errorf("error when extract jwt token")
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *UserInfo {
	raw, _ := ctx.Value(userCtxKey).(*UserInfo)
	return raw
}

func ForIpContext(ctx context.Context) string {
	raw, _ := ctx.Value(ipCtxKey).(string)
	return raw
}

func getUserInfoFromRedis(ctx context.Context, c *base.Client, apiToken string) *UserInfo {
	userInfoMap, err := c.RE.HGetAll(ctx, "tokens:"+apiToken).Result()
	var userInfo UserInfo
	if err == nil && len(userInfoMap) > 0 {
		userId, _ := strconv.Atoi(userInfoMap["user_id"])
		if userId == 0 {
			return nil
		}

		userInfo.UserId = int64(userId)
		userInfo.Role = userInfoMap["role"]
		ip := ForIpContext(ctx)
		c.RE.HSet(ctx, "tokens:"+apiToken, "ip", ip, "access_time", time.Now())
		return &userInfo
	} else {
		user, err := model.QueryUserByToken(c.DB, apiToken)
		if err == nil {
			err2 := c.RE.HSet(ctx, "tokens:"+apiToken, "user_id", user.ID, "role", user.Role).Err()
			if err2 != nil {
				log.Println("error: redis set")
			}

			userInfo.UserId = user.ID
			userInfo.Role = user.Role
			return &userInfo
		} else {
			err2 := c.RE.HSet(ctx, "tokens:"+apiToken, "user_id", 0).Err()
			if err2 == nil {
				_ = c.RE.Expire(ctx, "tokens:"+apiToken, 10*time.Minute).Err()
			} else {
				log.Println("error: redis set")
			}
		}
	}
	return nil
}

func ReplyUnauthorized(w http.ResponseWriter) {
	resp := make(map[string]interface{})
	resp["code"] = fmt.Sprint(http.StatusUnauthorized)
	resp["message"] = "Request Unauthorized, please use API token or login first."
	w.Header().Set("Content-Type", "application/json")
	respJson, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(respJson)
}

func ReplyTooMany(w http.ResponseWriter) {
	resp := make(map[string]interface{})
	resp["code"] = fmt.Sprint(http.StatusTooManyRequests)
	resp["message"] = `Too Many Requests, your can slow down your access and wait for quota to recovery and use API "GET /rate_limite" to check your access quota,`
	w.Header().Set("Content-Type", "application/json")
	respJson, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write(respJson)

}

func GenRand(l int) string {
	b := make([]byte, l)
	m, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b[:m])
}
