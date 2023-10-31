package user

import (
	"net/http"
	"strconv"

	"gitee.com/Kashimura/go-baka-control/db/mdb"
	g "gitee.com/Kashimura/go-baka-control/global"
	"github.com/gin-gonic/gin"
)

func InsertAttendanceRecord(ctx *gin.Context) {
	user := ctx.PostForm("user")
	time := ctx.PostForm("time")
	device := ctx.PostForm("device")
	// 空数据不处理
	if user == "" || time == "" || device == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// 数值转为整数
	userid, err := strconv.Atoi(user)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	// 查询是否有该用户
	if mdb.DB.CheckUserID(userid) != nil {
		ctx.JSON(http.StatusOK, g.ResponseError("未找到该用户"))
		return
	}

	if err := mdb.DB.InsertCheckInOut(userid, time, device); err != nil {
		ctx.JSON(http.StatusBadRequest, g.ResponseError("数据插入失败"))
		return
	}

	ctx.JSON(http.StatusOK, g.ResponseSuccess(nil))
}

func Group(group *gin.RouterGroup) {
	group.POST("/attendance", InsertAttendanceRecord)
}
