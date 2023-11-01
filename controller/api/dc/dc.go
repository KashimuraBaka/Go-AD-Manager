package dc

import (
	"gitee.com/Kashimura/go-baka-control/services/pshell"
	"github.com/gin-gonic/gin"
)

func GetUserInfo(ctx *gin.Context) {
	pshell.Shell.GetUsers("")
}

func Group(group *gin.RouterGroup) {
	group.GET("/info", GetUserInfo)
}
