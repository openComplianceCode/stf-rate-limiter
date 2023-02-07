package base

import "time"

const (
	MAX_IDEL_CONNS     = 20
	MAX_OPEN_CONNS     = 80
	CONN_MAX_LIFE_TIME = time.Hour
	IP_HEADER_KEY      = "X-Real-IP"
	GITEE_REDIRECT     = "/gitee_redirect"
	GITHUB_REDIRECT    = "/github_redirect"
	GITEE_API          = "https://gitee.com/api/v5/user"
	GITHUB_API         = "https://api.github.com/user"
	GITEE_OAUTH_CODE   = "https://gitee.com/oauth/authorize?response_type=code"
	GITEE_OAUTH_TOKEN  = "https://gitee.com/oauth/token?grant_type=authorization_code"
	GITHUB_OAUTH_CODE  = "https://github.com/login/oauth/authorize?access_type=online"
	GITHUB_OAUTH_TOKEN = "https://github.com/login/oauth/access_token"
	UpstreamServer     = "compliance.openeuler.org"

	CMD_DIR       = "/cmds/"
	BASE_TMP_DIR  = "/app/temp/"
	BASE_REPO_DIR = "/app/repos/"
)

var (
	GeneralRL = map[string]interface{}{
		"nologin":  1000,
		"everyone": 2000,
		"fans":     4000,
		"gold":     16000,
		"diamon":   64000,
		"root":     100000000,
	}
)
