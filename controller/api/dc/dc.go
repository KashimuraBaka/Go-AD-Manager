package dc

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	g "gitee.com/Kashimura/go-baka-control/global"
	"gitee.com/Kashimura/go-baka-control/services/pshell"
	"gitee.com/Kashimura/go-baka-control/services/tools"
	"github.com/gin-gonic/gin"
)

func GetUserInfo(ctx *gin.Context) {
	user := ctx.Query("user")
	if (user != "" && strings.ToLower(user) == "administrator") || tools.HasSysbol(user) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	users, err := pshell.Shell.GetUsers(user)
	if err != nil {
		ctx.JSON(http.StatusOK, g.ResponseError("查询失败"))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseSuccess(users, "查询成功"))
	}
}

func EnableUser(ctx *gin.Context) {
	user, unlockStr := ctx.PostForm("user"), ctx.PostForm("unlock")
	unlock, err := strconv.ParseBool(unlockStr)
	// 禁止管理员操作
	if err != nil || user == "" || strings.ToLower(user) == "administrator" || tools.HasSysbol(user) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// 解锁失败
	if err := pshell.Shell.EnableUser(user, unlock); err != nil {
		ctx.JSON(http.StatusOK, g.ResponseError(fmt.Sprintf("[%s] 修改失败!原因：可能未找到该账号", user)))
		return
	}
	// 返回结果
	if unlock {
		ctx.JSON(http.StatusOK, g.ResponseSuccess(nil, fmt.Sprintf("%s 解锁完毕!", user)))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseSuccess(nil, fmt.Sprintf("%s 锁定完毕!", user)))
	}
}

func UnlockUser(ctx *gin.Context) {
	user, password := ctx.PostForm("user"), ctx.PostForm("password")
	// 禁止管理员操作
	if user == "" || strings.ToLower(user) == "administrator" || tools.HasSysbol(user) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// 解锁账户
	if err := pshell.Shell.UnlockUser(user, password); err != nil {
		ctx.JSON(http.StatusOK, g.ResponseError(fmt.Sprintf("[%s] 修改失败!原因：可能未找到该账号或密码复杂不符合要求", user)))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseSuccess(nil, fmt.Sprintf("重置完毕! 密码为:[%s]", password)))
	}
}

func Group(group *gin.RouterGroup) {
	group.GET("/info", GetUserInfo)
	group.POST("/enable", EnableUser)
	group.POST("/unlock", UnlockUser)
}
