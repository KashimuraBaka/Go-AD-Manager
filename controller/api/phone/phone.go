package phone

import (
	"net/http"

	"gitee.com/Kashimura/go-baka-control/db/mysql"
	g "gitee.com/Kashimura/go-baka-control/global"
	"github.com/gin-gonic/gin"
)

func GetPhoneList(ctx *gin.Context) {
	res := &[]PhoneInfo{}
	mysql.DB.Find(res)
	ctx.JSON(http.StatusOK, g.ResponseSuccess(res))
}

func UpdatePhoneInfo(ctx *gin.Context) {
	info := &PhoneInfo{
		Name:     ctx.PostForm("name"),
		Phone:    ctx.PostForm("phone"),
		RecordIP: ctx.RemoteIP(),
	}
	if mysql.DB.Create(info).Error != nil {
		ctx.JSON(http.StatusOK, g.ResponseError("update error"))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseSuccess(nil))
	}
}

func Group(group *gin.RouterGroup) {
	group.GET("/list", GetPhoneList)
	group.POST("/update", UpdatePhoneInfo)
}
