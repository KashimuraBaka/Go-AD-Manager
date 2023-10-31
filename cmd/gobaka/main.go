package gobaka

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gitee.com/Kashimura/go-baka-control/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func CustomLoggerFormatter(params gin.LogFormatterParams) string {
	switch params.StatusCode {
	case http.StatusOK:
	case http.StatusNotModified:
	default:
		return fmt.Sprintf("[gin] %s | %s | %d | %s %s \n",
			params.TimeStamp.Format("2006-01-02 15:04:05"),
			params.ClientIP,
			params.StatusCode,
			params.Method,
			params.Path,
		)
	}
	return ""
}

func Router() *gin.Engine {
	engine := gin.New()

	gin.DisableConsoleColor()

	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{Formatter: CustomLoggerFormatter}))

	// 开始Gzip压缩
	engine.Use(gzip.Gzip(gzip.BestSpeed))

	// 允许跨域
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 允许加载 static 目录静态网站
	engine.Use(static.Serve("/", static.LocalFile("static/", false)))
	// 未找到目录则进行自动跳转
	engine.NoRoute(func(ctx *gin.Context) {
		accept := ctx.Request.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			if ctx.Request.URL.Path == "/" {
				content, err := os.ReadFile("static/index.html")
				if (err) != nil {
					ctx.String(http.StatusNotFound, "Not Found")
					return
				}
				ctx.Data(http.StatusOK, "text/html; charset=utf-8", content)
			} else {
				ctx.String(http.StatusNotFound, "Not Found")
			}
		}
	})

	// 设置接口
	api.Api(engine.Group("api"))

	return engine
}
