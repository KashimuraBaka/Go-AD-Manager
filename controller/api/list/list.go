package list

import (
	"net/http"

	"gitee.com/Kashimura/go-baka-control/db/mdb"
	"gitee.com/Kashimura/go-baka-control/db/mysql"
	g "gitee.com/Kashimura/go-baka-control/global"
	"gitee.com/Kashimura/go-baka-control/services/jwt"
	"github.com/gin-gonic/gin"
)

func GetDomainUsers(ctx *gin.Context) {
	res := &[]DomainUser{}
	mysql.DB.Find(res)
	ctx.JSON(http.StatusOK, g.ResponseSuccess(res, "success"))
}

func GetSystemUrls(ctx *gin.Context) {
	res := &[]SystemUrl{}
	mysql.DB.Find(res)
	ctx.JSON(http.StatusOK, g.ResponseSuccess(res, "success"))
}

func GetDownloadFiles(ctx *gin.Context) {
	res := []DownloadFile{}
	if mysql.DB.Find(&res).RowsAffected > 0 {
		for i, v := range res {
			res[i].FileName = jwt.GetFileToken(v.ID)
		}
	}
	ctx.JSON(http.StatusOK, g.ResponseSuccess(res, "success"))
}

func GetAttendanceList(ctx *gin.Context) {
	users, err := mdb.DB.SelectUserInfo()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, g.ResponseError("查询错误"))
		return
	}
	ctx.JSON(http.StatusOK, g.ResponseSuccess(users, "success"))
}

func Group(group *gin.RouterGroup) {
	group.GET("/domain", GetDomainUsers)
	group.GET("/urls", GetSystemUrls)
	group.GET("/downloads", GetDownloadFiles)
	group.GET("/attendance", GetAttendanceList)
}
