package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/alec-z/rate_limiter/auth"
	"github.com/alec-z/rate_limiter/base"
	"github.com/alec-z/rate_limiter/model"
	"github.com/alec-z/rate_limiter/util"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

func constsHandler(w http.ResponseWriter, r *http.Request) {
	consts := map[string]map[string]interface{}{
		"authURL": {
			"github": auth.GenerateURL("github"),
			"gitee":  auth.GenerateURL("gitee"),
		},
		"quota": base.GeneralRL,
	}
	util.ReplyJson(w, http.StatusOK, consts)
}

func showAPITokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if userInfo := auth.ForContext(ctx); userInfo != nil {
		if user, err := model.QueryUser(client.DB, userInfo.UserId); err == nil {
			var apiTokenPrefix string
			if user.APIToken == nil {
				apiTokenPrefix = ""
			} else if len(*user.APIToken) > 4 {
				apiTokenPrefix = (*user.APIToken)[:4]
			}
			resp := map[string]interface{}{
				"api_token_generate_time": user.APITokenGenerateTime,
				"api_token_prefix":        apiTokenPrefix,
			}
			if user.APITokenGenerateTime != nil {
				rm, err2 := client.RE.HGetAll(ctx, "token:"+*user.APIToken).Result()
				if err2 == nil {
					resp["ip"] = rm["ip"]
					resp["access_time"] = rm["access_time"]
				}
			}
			util.ReplyJson(w, http.StatusOK, resp)
			return
		}
	}
	auth.ReplyUnauthorized(w)
}
func generateAPITokenHandler(w http.ResponseWriter, r *http.Request) {
	var userInfo *auth.UserInfo
	ctx := r.Context()
	if userInfo = auth.ForContext(ctx); userInfo == nil {
		auth.ReplyUnauthorized(w)
		return
	}
	if user, err := model.QueryUser(client.DB, userInfo.UserId); err == nil {
		var claims auth.CustomClaims
		claims.UserId = user.ID
		claims.Role = string(user.Role)
		oldApiToken := *user.APIToken
		newApiToken := auth.GenRand(64)

		user.APIToken = &newApiToken
		now := time.Now()
		user.APITokenGenerateTime = &now
		model.UpdateUser(client.DB, user)
		client.RE.Del(ctx, "tokens:"+oldApiToken)

		resp := map[string]interface{}{
			"api_token": newApiToken,
			"message":   "if you have the old api_token, it has expired. Use this in 'Authorization' header",
		}
		util.ReplyJson(w, http.StatusOK, resp)
		return

	}
	respErr := map[string]interface{}{
		"api_token": "",
		"message":   "error when get App APi Token",
	}
	util.ReplyJson(w, http.StatusOK, respErr)
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	var userInfo *auth.UserInfo
	ctx := r.Context()
	if userInfo = auth.ForContext(ctx); userInfo == nil {
		auth.ReplyUnauthorized(w)
		return
	}
	user, err := model.QueryUser(client.DB, userInfo.UserId)
	if user != nil && err == nil {
		util.ReplyJson(w, http.StatusOK, user)

	}

}

func updateUserDetailHandler(w http.ResponseWriter, r *http.Request) {
	var userDetail model.User
	err := json.NewDecoder(r.Body).Decode(&userDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if userInfo := auth.ForContext(ctx); userInfo == nil {
		auth.ReplyUnauthorized(w)
		return
	} else {
		userDetail.ID = userInfo.UserId
	}

	if userDetail.FirstName == nil || *userDetail.FirstName == "" {
		http.Error(w, "Must have first name", http.StatusBadRequest)
		return
	}
	if userDetail.LastName == nil || *userDetail.LastName == "" {
		http.Error(w, "must have last name", http.StatusBadRequest)
		return
	}
	if userDetail.EmailAddress == nil || *userDetail.EmailAddress == "" {
		http.Error(w, "must have email address", http.StatusBadRequest)
		return
	}
	if userDetail.City == nil || *userDetail.City == "" {
		http.Error(w, "must have city", http.StatusBadRequest)
		return
	}
	if userDetail.Country == nil || *userDetail.Country == "" {
		http.Error(w, "must have country", http.StatusBadRequest)
		return
	}
	if userDetail.PostalCode == nil || *userDetail.PostalCode == "" {
		http.Error(w, "must have postal code", http.StatusBadRequest)
		return
	}

	if userDetail.Address == nil || *userDetail.Address == "" {
		http.Error(w, "must have address", http.StatusBadRequest)
		return
	}

	newUser, err := model.UpdateUserDetail(client.DB, &userDetail)
	if err != nil {
		http.Error(w, "update user's profile detail failed", http.StatusInternalServerError)
		return
	}
	util.ReplyJson(w, http.StatusOK, newUser)
}

func authRedirectHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successfully!"))
}

func checkRemainingHandler(w http.ResponseWriter, r *http.Request) {
	pass, remaining := auth.CheckRateLimit(r.Context(), client.RE, 1)
	resp := map[string]interface{}{
		"more":      pass,
		"remaining": remaining,
		"message":   "check the remaining will also comsume 1 quota",
	}
	util.ReplyJson(w, http.StatusOK, resp)
}

var client base.Client

func main() {
	setupClient(&client)

	githubOauth2 := auth.Oauth2{T: "github", Client: &client}
	giteeOauth2 := auth.Oauth2{T: "gitee", Client: &client}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/consts", constsHandler)
	mux.Handle(base.GITEE_REDIRECT, auth.Middleware(&client, &giteeOauth2))
	mux.Handle(base.GITHUB_REDIRECT, auth.Middleware(&client, &githubOauth2))

	mux.Handle("/api/me", auth.Middleware(&client, http.HandlerFunc(meHandler)))
	mux.Handle("/api/user_detail", auth.Middleware(&client, http.HandlerFunc(updateUserDetailHandler)))

	mux.Handle("/api/generate_api_token", auth.Middleware(&client, http.HandlerFunc(generateAPITokenHandler)))
	mux.Handle("/api/show_api_token", auth.Middleware(&client, http.HandlerFunc(showAPITokenHandler)))

	mux.Handle("/auth-redirect", auth.Middleware(&client, http.HandlerFunc(authRedirectHandler)))
	mux.Handle("/api/quota_remaining", auth.Middleware(&client, http.HandlerFunc(checkRemainingHandler)))

	mux.Handle("/api/", auth.Middleware(&client, http.HandlerFunc(getForwardHandler())))

	log.Println("server setup at localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}

func setupClient(client *base.Client) {
	dbClient, err := OpenDB()
	if err != nil {
		log.Fatalln("cannot connect to DB !", err)
		return
	}
	reClient, err := OpenRedis()
	if err != nil {
		log.Fatalln("cannot connect to redis !", err)
		return
	}
	client.DB = dbClient
	client.RE = reClient

}

func OpenDB() (*sql.DB, error) {
	mysqlDataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/rate_limiter?parseTime=true", base.MysqlUser,
		base.MysqlPassword, base.MysqlHost, base.MysqlPort)

	drv, err := sql.Open("mysql", mysqlDataSource)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Get the underlying sql.DB object of the driver.
	drv.SetMaxIdleConns(base.MAX_IDEL_CONNS)
	drv.SetMaxOpenConns(base.MAX_OPEN_CONNS)
	drv.SetConnMaxLifetime(base.CONN_MAX_LIFE_TIME)
	return drv, nil
}

func OpenRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", base.RedisHost, base.RedisPort),
		Password: base.RedisPassword,
	})
	if rdb != nil {
		return rdb, nil
	} else {
		return nil, fmt.Errorf("cannot connect to redis")
	}
}

func getForwardHandler() func(http.ResponseWriter, *http.Request) {
	fUrl, _ := url.Parse("https://" + base.UpstreamServer)
	proxy := httputil.NewSingleHostReverseProxy(fUrl)
	pathConfig, _ := util.ReadPathConfig()
	log.Println(pathConfig)
	return func(w http.ResponseWriter, r *http.Request) {
		if r != nil {
			r.Host = base.UpstreamServer
			var consume int = 1
			for _, p := range pathConfig.Paths {
				if p.PathType == "Exact" && r.URL.Path == p.Path {
					consume = p.Consumption
					break
				} else if p.PathType == "Prefix" && strings.HasPrefix(r.URL.Path, p.Path) {
					consume = p.Consumption
					break
				}
			}
			pass, _ := auth.CheckRateLimit(r.Context(), client.RE, consume)
			if !pass {
				http.Error(w, "Too many request, exceed rate limit", http.StatusTooManyRequests)
				return
			}
			proxy.ServeHTTP(w, r)

		}
	}
}
