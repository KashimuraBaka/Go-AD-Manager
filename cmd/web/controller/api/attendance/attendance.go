package attendance

import (
	"net/http"
	"strconv"

	"gitee.com/Kashimura/go-baka-control/cmd/web/controller/middleware"
	"gitee.com/Kashimura/go-baka-control/services/db/mdb"
	"gitee.com/Kashimura/go-baka-control/utils/webhttp"
	"github.com/gin-gonic/gin"
)

func Register(parentGroup *gin.RouterGroup) {
	group := parentGroup.Group("attendance")

	group.Use(middleware.VerifyAccountToken())
	group.GET("/list", GetAttendanceList)
	group.POST("/insert", InsertAttendanceRecord)
}

func GetAttendanceList(ctx *gin.Context) {
	users, err := mdb.DB.SelectUserInfo()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "查询错误: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Data: users,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}

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
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "未找到该用户",
		})
		return
	}

	if err := mdb.DB.InsertCheckInOut(userid, time, device); err != nil {
		ctx.JSON(http.StatusBadRequest, webhttp.Response{
			Code: webhttp.ERROR_DATABASE_CREATE_ERROR,
			Msg:  "数据插入失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}
