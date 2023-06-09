package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/alec-z/rate_limiter/base"
	"github.com/alec-z/rate_limiter/model"
)

type Oauth2 struct {
	T string //github or github
	*base.Client
}

func (o *Oauth2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ui := ForContext(r.Context()); ui != nil && r.URL.Query().Get("state") == "" {
		http.Redirect(w, r, "/auth-redirect", http.StatusFound)
		return
	}

	var u *model.User

	if ui := ForContext(r.Context()); ui != nil && r.URL.Query().Get("state") != "" {
		u, _ = model.QueryUser(o.Client.DB, ui.UserId) // already login and refresh
	}

	code := r.URL.Query().Get("code")
	log.Println("code is : ", code)

	accessToken, err := exchangeAccessToken(r.Context(), o.T, code, r.URL.Query().Get("state"))
	if err != nil {
		log.Println("exchange token failed: ", err)
		return
	}

	uj, err := requestForUserJson(r.Context(), o.T, accessToken)
	if err != nil {
		log.Println("request for user json failed: ", err)
		return
	}

	userID := fmt.Sprint(uj.ID)

	//query and update object
	if u == nil {
		if o.T == "gitee" {
			u, err = model.QueryUserByGiteeID(o.DB, userID)
		} else if o.T == "github" {
			u, err = model.QueryUserByGithubID(o.DB, userID)
		}
	}
	if err != nil {
		var userMutation model.User
		setUserInfo(&userMutation, o.T, uj)
		u, err = model.CreateUser(o.DB, &userMutation)
	} else {
		setUserInfo(u, o.T, uj)
		u, err = model.UpdateUser(o.DB, u)
	}
	if err != nil {
		log.Printf("create or update user info error", err)
		return
	}

	//createJwt
	var c CustomClaims
	c.UserId = u.ID
	c.Role = string(u.Role)
	if tokenString, err2 := CreateOrRefreshToken(&c); err2 == nil {
		cookie := http.Cookie{Name: "auth-cookie", Value: tokenString, MaxAge: 60 * (JWT_VALID_MINS - 3)}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/auth-redirect", http.StatusFound)
		return
	}

	http.Error(w, "Invalid oauth2 authorization", http.StatusUnauthorized)
}

func setUserInfo(m *model.User, authType string, uj *userJson) {
	sID := fmt.Sprint(uj.ID)
	if authType == "gitee" {
		m.GiteeID = &sID
		m.GiteeEmail = &uj.Email
		m.GiteeLogin = &uj.Login
		m.GiteeAvatarUrl = &uj.AvatarURL
		m.GiteeName = &uj.Name
	} else if authType == "github" {
		m.GithubID = &sID
		m.GithubEmail = &uj.Email
		m.GithubLogin = &uj.Login
		m.GithubAvatarUrl = &uj.AvatarURL
		m.GithubName = &uj.Name
	}
}

func exchangeAccessToken(ctx context.Context, authType string, code string, state string) (accessToken string, err error) {

	var values url.Values
	var rUrl string

	if authType == "github" {
		rUrl = base.GITHUB_OAUTH_TOKEN
		values = url.Values{
			"client_id":     {base.GithubClientId},
			"client_secret": {base.GithubClientSecret},
			"code":          {code},
		}
	} else if authType == "gitee" {
		rUrl = base.GITEE_OAUTH_TOKEN
		redirectUrl := base.SelfDomain + base.GITEE_REDIRECT
		if state != "" {
			redirectUrl += "?state=" + state
		}
		values = url.Values{
			"client_id":     {base.GiteeClientId},
			"client_secret": {base.GiteeClientSecret},
			"code":          {code},
			"redirect_uri":  {redirectUrl},
			"grant_type":    {"authorization_code"},
		}
	}
	log.Println(authType)
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", rUrl, strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	var token map[string]interface{}
	if tokenBytes, err := io.ReadAll(resp.Body); err == nil {
		if err = json.Unmarshal([]byte(tokenBytes), &token); err == nil {
			log.Println("access token", token)
			if accessToken = token["access_token"].(string); accessToken != "" {
				return accessToken, nil
			} else {
				return "", fmt.Errorf("error access token")
			}
		} else {
			log.Println(err)
			return "", err
		}
	}
	return "", err
}

type userJson struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

func requestForUserJson(ctx context.Context, authType string, token string) (*userJson, error) {
	var uj userJson
	client := &http.Client{}
	var url string
	if authType == "github" {
		url = base.GITHUB_API
	} else if authType == "gitee" {
		url = base.GITEE_API
	}

	if req, err := http.NewRequestWithContext(ctx, "GET", url, nil); err == nil {
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Accept", "application/json")
		if resp, err := client.Do(req); err == nil {
			if respBytes, err := io.ReadAll(resp.Body); err == nil {
				if err = json.Unmarshal(respBytes, &uj); err == nil {
					log.Println(string(respBytes))
					log.Println(uj)
					return &uj, nil
				}
			} else {
				log.Println(err)
				return nil, err
			}
		}
	}
	return nil, fmt.Errorf("request user info fail")
}

func GenerateURL(authType string) string {
	var d, s, r string
	if authType == "github" {
		d = base.GITHUB_OAUTH_CODE
		s = base.GithubClientId
		r = base.GITHUB_REDIRECT
	} else if authType == "gitee" {
		d = base.GITEE_OAUTH_CODE
		s = base.GiteeClientId
		r = base.GITEE_REDIRECT
	}
	return fmt.Sprintf("%s&client_id=%s&redirect_uri=%s", d, s, base.SelfDomain+r)

}
