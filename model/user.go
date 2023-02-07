package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID                   int64      `json:"id"`
	CreateTime           time.Time  `json:"createTime"`
	UpdateTime           time.Time  `json:"updateTime"`
	Role                 string     `json:"role"`
	GiteeID              *string    `json:"giteeID"`
	GiteeLogin           *string    `json:"giteeLogin"`
	GiteeName            *string    `json:"giteeLogin"`
	GiteeEmail           *string    `json:"giteeEmail"`
	GiteeAvatarUrl       *string    `json:"giteeAvatarUrl"`
	GithubID             *string    `json:"githubID"`
	GithubLogin          *string    `json:"githubLogin"`
	GithubName           *string    `json:"githubName"`
	GithubEmail          *string    `json:"githubEmail"`
	GithubAvatarUrl      *string    `json:"githubAvatarUrl"`
	APIToken             *string    `json:"apiToken"`
	APITokenGenerateTime *time.Time `json:"apiTokenGenerateTime"`
	UserDetail
}
type UserDetail struct {
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	EmailAddress *string `json:"emailAddress"`
	Company      *string `json:"company"`
	City         *string `json:"city"`
	Country      *string `json:"country"`
	PostalCode   *string `json:"postalCode"`
	Address      *string `json:"address"`
	AboutMe      *string `json:"aboutMe"`
}

func CreateUser(db *sql.DB, user *User) (*User, error) {
	res, err := db.Exec(`insert into users(create_time, update_time, role, 
		gitee_id, gitee_login, gitee_name, gitee_email, gitee_avatar_url,
		github_id, github_login, github_name, github_email, github_avatar_url) values(Now(), Now(), "everyone", ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.GiteeID, user.GiteeLogin, user.GiteeName, user.GiteeEmail, user.GiteeAvatarUrl,
		user.GithubID, user.GithubLogin, user.GithubName, user.GithubEmail, user.GithubAvatarUrl)
	if err != nil {
		return nil, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return QueryUser(db, lastId)
}

func UpdateUser(db *sql.DB, user *User) (*User, error) {

	oldUser, err := QueryUser(db, user.ID)
	if err != nil {
		return nil, err
	}
	if oldUser.GiteeID != nil && user.GiteeID == nil {
		user.GiteeID = oldUser.GiteeID
	}
	if oldUser.GiteeLogin != nil && user.GiteeLogin == nil {
		user.GiteeLogin = oldUser.GiteeLogin
	}
	if oldUser.GiteeName != nil && user.GiteeName == nil {
		user.GiteeName = oldUser.GiteeName
	}
	if oldUser.GiteeEmail != nil && user.GiteeEmail == nil {
		user.GiteeEmail = oldUser.GiteeEmail
	}
	if oldUser.GiteeAvatarUrl != nil && user.GiteeAvatarUrl == nil {
		user.GiteeAvatarUrl = oldUser.GiteeAvatarUrl
	}
	if oldUser.GithubID != nil && user.GithubID == nil {
		user.GithubID = oldUser.GithubID
	}
	if oldUser.GithubLogin != nil && user.GithubLogin == nil {
		user.GithubLogin = oldUser.GithubLogin
	}
	if oldUser.GithubName != nil && user.GithubName == nil {
		user.GithubName = oldUser.GithubName
	}
	if oldUser.GithubEmail != nil && user.GithubEmail == nil {
		user.GithubEmail = oldUser.GithubEmail
	}
	if oldUser.GithubAvatarUrl != nil && user.GithubAvatarUrl == nil {
		user.GithubAvatarUrl = oldUser.GithubAvatarUrl
	}

	if oldUser.Role != "" && user.Role == "" {
		user.Role = oldUser.Role
	}
	if oldUser.APIToken != nil && user.APIToken == nil {
		user.APIToken = oldUser.APIToken
	}

	if oldUser.APITokenGenerateTime != nil && user.APITokenGenerateTime == nil {
		user.APITokenGenerateTime = oldUser.APITokenGenerateTime
	}

	res, err := db.Exec(`update users set update_time = Now(), role = ?, 
		gitee_id = ?, gitee_login = ?, gitee_name = ?, gitee_email = ?, gitee_avatar_url = ?,
		github_id = ?, github_login = ?, github_name = ?, github_email = ?, github_avatar_url = ?, api_token = ?, api_token_generate_time = ? where id = ?`,
		user.Role,
		user.GiteeID, user.GiteeLogin, user.GiteeName, user.GiteeEmail, user.GiteeAvatarUrl,
		user.GithubID, user.GithubLogin, user.GithubName, user.GithubEmail, user.GithubAvatarUrl, user.APIToken, user.APITokenGenerateTime, user.ID)
	if err != nil {
		return nil, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserDetail(db *sql.DB, user *User) (*User, error) {

	oldUser, err := QueryUser(db, user.ID)
	if err != nil {
		return nil, err
	}
	if oldUser.FirstName != nil && user.FirstName == nil {
		user.FirstName = oldUser.FirstName
	}
	if oldUser.LastName != nil && user.LastName == nil {
		user.LastName = oldUser.LastName
	}
	if oldUser.EmailAddress != nil && user.EmailAddress == nil {
		user.EmailAddress = oldUser.EmailAddress
	}
	if oldUser.Company != nil && user.Company == nil {
		user.Company = oldUser.Company
	}
	if oldUser.City != nil && user.City == nil {
		user.City = oldUser.City
	}
	if oldUser.Country != nil && user.Country == nil {
		user.Country = oldUser.Country
	}
	if oldUser.PostalCode != nil && user.PostalCode == nil {
		user.PostalCode = oldUser.PostalCode
	}
	if oldUser.Address != nil && user.Address == nil {
		user.Address = oldUser.Address
	}
	if oldUser.AboutMe != nil && user.AboutMe == nil {
		user.AboutMe = oldUser.AboutMe
	}

	res, err := db.Exec(`update users set update_time = Now(), first_name = ?, 
		last_name = ?, email_address = ?, company = ?, city = ?, country = ?,
		postal_code = ?, address = ?, about_me = ? where id = ?`,
		user.FirstName,
		user.LastName, user.EmailAddress, user.Company, user.City, user.Country,
		user.PostalCode, user.Address, user.AboutMe, user.ID)
	if err != nil {
		return nil, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if newUser, err2 := QueryUser(db, user.ID); err2 == nil {
		return newUser, nil
	} else {
		return nil, err2
	}

}

func QueryUser(db *sql.DB, id int64) (*User, error) {
	var user User
	err := db.QueryRow(`select id, create_time, update_time, role, 
	gitee_id, gitee_login, gitee_name, gitee_email, gitee_avatar_url,
	github_id, github_login, github_name, github_email, github_avatar_url, api_token, api_token_generate_time,
	first_name, last_name, email_address, company, city, country, postal_code, address, about_me
	 from users where id = ?`, id).Scan(
		&user.ID,
		&user.CreateTime,
		&user.UpdateTime,
		&user.Role,
		&user.GiteeID,
		&user.GiteeLogin,
		&user.GiteeName,
		&user.GiteeEmail,
		&user.GiteeAvatarUrl,
		&user.GithubID,
		&user.GithubLogin,
		&user.GithubName,
		&user.GithubEmail,
		&user.GithubAvatarUrl,
		&user.APIToken,
		&user.APITokenGenerateTime,
		&user.FirstName,
		&user.LastName,
		&user.EmailAddress,
		&user.Company,
		&user.City,
		&user.Country,
		&user.PostalCode,
		&user.Address,
		&user.AboutMe,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func QueryUserByToken(db *sql.DB, apiToken string) (*User, error) {
	var user User
	err := db.QueryRow(`select id, role, api_token, api_token_generate_time from users
	 where api_token = ?`, apiToken).Scan(
		&user.ID,
		&user.Role,
		&user.APIToken,
		&user.APITokenGenerateTime,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func QueryUserByGiteeID(db *sql.DB, giteeID string) (*User, error) {
	var user User
	err := db.QueryRow(`select id, create_time, update_time, role, 
	gitee_id, gitee_login, gitee_name, gitee_email, gitee_avatar_url,
	github_id, github_login, github_name, github_email, github_avatar_url  from users where gitee_id = ?`, giteeID).Scan(
		&user.ID,
		&user.CreateTime,
		&user.UpdateTime,
		&user.Role,
		&user.GiteeID,
		&user.GiteeLogin,
		&user.GiteeName,
		&user.GiteeEmail,
		&user.GiteeAvatarUrl,
		&user.GithubID,
		&user.GithubLogin,
		&user.GithubName,
		&user.GithubEmail,
		&user.GithubAvatarUrl,
		&user.APIToken)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func QueryUserByGithubID(db *sql.DB, githubID string) (*User, error) {
	var user User
	err := db.QueryRow(`select id, create_time, update_time, role, 
	gitee_id, gitee_login, gitee_name, gitee_email, gitee_avatar_url,
	github_id, github_login, github_name, github_email, github_avatar_url, api_token from users where github_id = ?`, githubID).Scan(
		&user.ID,
		&user.CreateTime,
		&user.UpdateTime,
		&user.Role,
		&user.GiteeID,
		&user.GiteeLogin,
		&user.GiteeName,
		&user.GiteeEmail,
		&user.GiteeAvatarUrl,
		&user.GithubID,
		&user.GithubLogin,
		&user.GithubName,
		&user.GithubEmail,
		&user.GithubAvatarUrl,
		&user.APIToken)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
