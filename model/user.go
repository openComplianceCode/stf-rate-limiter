package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID              int64
	CreateTime      time.Time
	UpdateTime      time.Time
	Role            string
	GiteeID         *string
	GiteeLogin      *string
	GiteeName       *string
	GiteeEmail      *string
	GiteeAvatarUrl  *string
	GithubID        *string
	GithubLogin     *string
	GithubName      *string
	GithubEmail     *string
	GithubAvatarUrl *string
	APIKey          string
}

func CreateUser(db *sql.DB, user *User) (*User, error) {
	res, err := db.Exec(`insert into users(create_time, update_time, role, 
		gitee_id, gitee_login, gitee_name, gitee_email, gitee_avatar_url,
		github_id, github_login, github_name, github_email, github_avatar_url, api_key) values(Now(), Now(), "everyone", ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
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
	if oldUser.APIKey != "" && user.APIKey == "" {
		user.Role = oldUser.Role
	}

	res, err := db.Exec(`update users set update_time = ? , role = ?, 
		gitee_id = ?, gitee_login = ?, gitee_name = ?, gitee_email = ?, gitee_avatar_url = ?,
		github_id = ?, github_login = ?, github_name = ?, github_email = ?, github_avatar_url = ?, api_key = ?) where id = ?`,
		user.GiteeID, user.GiteeLogin, user.GiteeName, user.GiteeEmail, user.GiteeAvatarUrl,
		user.GithubID, user.GithubLogin, user.GithubName, user.GithubEmail, user.GithubAvatarUrl, user.APIKey, user.ID)
	if err != nil {
		return nil, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func QueryUser(db *sql.DB, id int64) (*User, error) {
	var user User
	err := db.QueryRow(`select id, create_time, update_time, role, 
	gitee_id, gitee_login, gitee_name, gitee_email, gitee_avatar_url,
	github_id, github_login, github_name, github_email, github_avatar_url, api_key from users where id = ?`, id).Scan(
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
		&user.APIKey)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
