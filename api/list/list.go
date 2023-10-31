package list

import (
	"net/http"

	"gitee.com/Kashimura/go-baka-control/db/mdb"
	"gitee.com/Kashimura/go-baka-control/db/mysql"
	g "gitee.com/Kashimura/go-baka-control/global"
	"github.com/gin-gonic/gin"
)

func GetDomainUsers(ctx *gin.Context) {
	res := &[]DomainUser{}
	mysql.DB.Find(res)
	ctx.JSON(http.StatusOK, g.ResponseSuccess(res))
}

func GetSystemUrls(ctx *gin.Context) {
	res := &[]SystemUrl{}
	mysql.DB.Find(res)
	ctx.JSON(http.StatusOK, g.ResponseSuccess(res))
}

func GetDownloadFiles(ctx *gin.Context) {
	res := &[]DownloadFile{}
	mysql.DB.Find(res)
	ctx.JSON(http.StatusOK, g.ResponseSuccess(res))
}

func GetAttendanceList(ctx *gin.Context) {
	users, err := mdb.DB.SelectUserInfo()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, g.ResponseError("查询错误"))
		return
	}
	ctx.JSON(http.StatusOK, g.ResponseSuccess(users))
}

func Group(group *gin.RouterGroup) {
	group.GET("/domain", GetDomainUsers)
	group.GET("/urls", GetSystemUrls)
	group.GET("/downloads", GetDownloadFiles)
	group.GET("/attendance", GetAttendanceList)
}
