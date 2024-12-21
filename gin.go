package xlogger

import (
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// GinLogConfigure http server log settings.
type GinLogConfigure struct {
	// skip url and do not log
	skipPrefix []string

	// enable request body
	requestBody bool

	// ip log ip
	ipFunc func(ctx *gin.Context) string

	// user-agent
	userAgent bool

	// headerKeys
	headerKeys []string
}

func (c *GinLogConfigure) SkipPrefix(prefix ...string) {
	c.skipPrefix = append(c.skipPrefix, prefix...)
}

func (c *GinLogConfigure) EnableRequestBody() {
	c.requestBody = true
}

func (c *GinLogConfigure) LogIP(ipFunc func(ctx *gin.Context) string) {
	c.ipFunc = ipFunc
}

func (c *GinLogConfigure) EnableUserAgent() {
	c.userAgent = true
}

func (c *GinLogConfigure) AddHeaderKeys(key ...string) {
	c.headerKeys = append(c.headerKeys, key...)
}

// GinLog http server Middleware Logger
func GinLog(conf GinLogConfigure) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		for _, v := range conf.skipPrefix {
			if strings.HasPrefix(c.Request.RequestURI, v) {
				return
			}
		}
		var field = make(Field)
		// 结束时间
		end := time.Now()
		// 执行时间
		rt := end.Sub(start) / time.Millisecond
		field["duration"] = rt
		// 状态码
		code := c.Writer.Status()
		field["status"] = code
		if conf.ipFunc != nil {
			ip := conf.ipFunc(c)
			field["ip"] = ip
		}
		for _, header := range conf.headerKeys {
			val := c.Request.Header.Get(header)
			field["header_"+header] = val
		}
		if c.Err() != nil {
			field["error"] = c.Err()
		}
		field["method"] = c.Request.Method
		field["uri"] = c.Request.URL.String()
		if conf.userAgent {
			field["userAgent"] = c.Request.UserAgent()
		}
		field["scene"] = "http_server_request"
		WithContext(c.Request.Context()).WithFields(field).Info("http_request")
	}
}
