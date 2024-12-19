package xlogger

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGinLogger(t *testing.T) {
	g := gin.New()
	g.Use(GinLog(GinLogConfigure{
		requestBody: true,
	}))

	g.POST("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"foo": "bar",
		})
	})

	go func() {
		g.Run(":9000")
	}()

	time.Sleep(1 * time.Second)

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:9000", strings.NewReader(`{"foo":"bar"}`))
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
}
