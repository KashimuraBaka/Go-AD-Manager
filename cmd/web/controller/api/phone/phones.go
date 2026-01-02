package list

import (
	"net/http"

	"gitee.com/Kashimura/go-baka-control/services/db/mysql"
	"gitee.com/Kashimura/go-baka-control/utils/webhttp"
	"github.com/gin-gonic/gin"
)

func Register(parentGruop *gin.RouterGroup) {
	group := parentGruop.Group("/phone")

	group.GET("/list", GetPhoneList)
	group.POST("/update", UpdatePhoneInfo)
}

func GetPhoneList(ctx *gin.Context) {
	res := &[]mysql.PhoneInfo{}
	mysql.DB.Find(res)
	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Data: res,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}

func UpdatePhoneInfo(ctx *gin.Context) {
	info := &mysql.PhoneInfo{
		Name:     ctx.PostForm("name"),
		Phone:    ctx.PostForm("phone"),
		RecordIP: ctx.RemoteIP(),
	}
	if mysql.DB.Create(info).Error != nil {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_UPLOAD_ERROR,
			Msg:  webhttp.GetMessageByCode(webhttp.ERROR_UPLOAD_ERROR),
		})
	} else {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
		})
	}
}
