package base

import (
	"os"
)

var (
	GiteeClientId      string = "3680e4d97dc506dfdc50c8c11588f0fe703e64367bacf0f33738c392889fdc09"
	GiteeClientSecret  string = ""
	GithubClientId     string = "9bebf3d45dc11bbc454a"
	GithubClientSecret string = ""

	MysqlUser     string = "root"
	MysqlPassword        = "root"
	MysqlHost            = "127.0.0.1"
	MysqlPort            = "3306"

	RedisPassword string = ""
	RedisHost            = "127.0.0.1"
	RedisPort            = "6379"

	ServerPort = "8080"

	HmacSecret = []byte("PTkCA7pu6maDQgf5uA7y/6NTqN3PSyXUJLLhKvLh0MRo")

	SelfDomain = "http://localhost:8080"
)

func init() {
	GiteeClientSecret = os.Getenv("GITEE_CLIENT_SECRET")
	GithubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	if env := os.Getenv("MYSQL_USER"); env != "" {
		MysqlUser = env
	}
	if env := os.Getenv("MYSQL_PASSWORD"); env != "" {
		MysqlPassword = env
	}
	if env := os.Getenv("MYSQL_HOST"); env != "" {
		MysqlHost = env
	}
	if env := os.Getenv("MYSQL_PORT"); env != "" {
		MysqlPort = env
	}

	if env := os.Getenv("REDIS_PASSWORD"); env != "" {
		RedisPassword = env
	}
	if env := os.Getenv("REDIS_HOST"); env != "" {
		RedisHost = env
	}
	if env := os.Getenv("REDIS_PORT"); env != "" {
		RedisPort = env
	}

	if env := os.Getenv("SEVER_PORT"); env != "" {
		ServerPort = env
	}

	if env := os.Getenv("HMAC_SECRET"); env != "" {
		HmacSecret = []byte(env)
	}

	if env := os.Getenv("SELF_DOMAIN"); env != "" {
		SelfDomain = env
	}

}
