package list

import (
	"net/http"

	"gitee.com/Kashimura/go-baka-control/services/db/mysql"
	"gitee.com/Kashimura/go-baka-control/utils/webhttp"
	"github.com/gin-gonic/gin"
)

func Register(parentGroup *gin.RouterGroup) {
	group := parentGroup.Group("/list")

	group.GET("/domain", GetDomainUsers)
	group.GET("/urls", GetSystemUrls)
}

func GetDomainUsers(ctx *gin.Context) {
	res := &[]mysql.DomainUser{}
	mysql.DB.Find(res)
	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}

func GetSystemUrls(ctx *gin.Context) {
	res := &[]mysql.SystemUrl{}
	mysql.DB.Find(res)
	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}
