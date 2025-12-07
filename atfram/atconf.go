package atfram

import (
	"maps"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/Dr-Deep/logging-go"
)

var (
	defaultHttpHeaders = map[string]string{
		"Accept-API-Version": "protocol=1.0,resource=2.0",                                                        // maybe that does something?
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:145.0) Gecko/20100101 Firefox/145.0", // js to be safe
		"Content-Type":       "application/json",
	}
)

type AldiTalkConfig struct {
	BaseURI string
	Headers map[string]string

	Username string
	Password string
}

func New(cfg *AldiTalkConfig, logger *logging.Logger) (*Client, error) {
	if cfg.BaseURI == "" {
		return nil, ErrAldiTalkConfigMissingBaseURI
	}

	if cfg.Username == "" {
		return nil, ErrAldiTalkConfigMissingUsername
	}

	if cfg.Password == "" {
		return nil, ErrAldiTalkConfigMissingPassword
	}

	// AldiTalk Config: join http headers
	var headers = defaultHttpHeaders
	maps.Copy(headers, cfg.Headers)
	cfg.Headers = headers

	// Http Client
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	var httpclient = &http.Client{
		Jar:     jar,
		Timeout: time.Millisecond * 3000,
	}

	// AldiTalk Client
	var c = &Client{
		c:      httpclient,
		logger: logger,
		atconf: cfg,
	}

	return c, nil
}
