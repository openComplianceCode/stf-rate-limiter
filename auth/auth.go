package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
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
	UserId int    `json:"user_id,omitempty"`
	Role   string `json:"role,omitempty"`
	ApiKey string `json:"api_key,omitempty"`
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
func Middleware(c *base.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			saveIPContext(r)
			var rr *http.Request = r

			if tokenString := strings.TrimSpace(r.Header.Get("Authorization")); tokenString != "" {
				rr = checkAuthorizationForLogin(tokenString, c, w, r)
			} else if cookie, err := r.Cookie("auth-cookie"); err == nil && c != nil && strings.TrimSpace(cookie.Value) != "" {
				rr = checkCookieForLogin(cookie, w, r)
			}

			if rr == nil {
				return
			}

			if !CheckRateLimit(rr.Context(), c.RE, 1) {
				http.Error(w, "Too many request, exceed rate limit", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, rr)
		})
	}
}

func CreateOrRefreshToken(c *CustomClaims) (tokenString string, err error) {
	c.ExpiresAt = time.Now().Add(time.Minute * JWT_VALID_MINS).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString(base.HmacSecret)
	return
}

func checkAuthorizationForLogin(tokenString string, c *base.Client, w http.ResponseWriter, r *http.Request) *http.Request {
	tokenStrArr := strings.Split(tokenString, " ")
	if len(tokenStrArr) > 1 {
		tokenString = tokenStrArr[len(tokenStrArr)-1]
	} else {
		tokenString = tokenStrArr[0]
	}
	if claims, err := extractClaims(tokenString); err == nil && claims.ApiKey != "" && claims.UserId != 0 {
		if apiKey := getApiKeyFromRedis(r.Context(), c, claims.UserId); apiKey == claims.ApiKey {
			return saveUserContext(claims, r)
		}
	}
	http.Error(w, "Invalid Authorization Token", http.StatusUnauthorized)
	return nil
}

func checkCookieForLogin(cookie *http.Cookie, w http.ResponseWriter, r *http.Request) *http.Request {
	tokenString := strings.TrimSpace(cookie.Value)
	if claims, err := extractClaims(tokenString); err == nil && claims.UserId != 0 { // check jwt token
		rr := saveUserContext(claims, r)
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

func saveUserContext(claims *CustomClaims, r *http.Request) *http.Request {
	if claims != nil {
		ctx := context.WithValue(r.Context(), userCtxKey, &claims.UserInfo)
		return r.WithContext(ctx)
	}
	return r
}

func saveIPContext(r *http.Request) *http.Request {
	remoteIP := r.Header.Get(base.IP_HEADER_KEY)
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

func getApiKeyFromRedis(ctx context.Context, c *base.Client, userId int64) string {
	if apiKey, err := c.RE.HGet(ctx, "api_keys", fmt.Sprint(userId)).Result(); err == nil {
		if user, err := model.QueryUser(c.DB, userId); err == nil {
			c.RE.HSet(ctx, "api_keys", fmt.Sprint(userId), apiKey)
			if user.APIKey != "" {
				return apiKey
			}
		}
	}
	return ""
}

func ReplyUnauthorized(ctx context.Context) {

}

func ReplyForbidden(ctx context.Context) {

}

func ReplyTooMany(ctx context.Context) {

}

func GenRand(l int) string {
	b := make([]byte, l)
	m, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b[:m])
}
