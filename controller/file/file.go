package file

import (
	"net/http"
	"net/url"
	"os"
	"path"

	"gitee.com/Kashimura/go-baka-control/db/mysql"
	"gitee.com/Kashimura/go-baka-control/services/jwt"
	"github.com/gin-gonic/gin"
)

func DownloadFile(ctx *gin.Context) {
	// 获取文件ID
	fileid, err := jwt.VerifyFileToken(ctx.Param("key"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	// 查询记录
	info := FileInfo{ID: fileid}
	if mysql.DB.Find(&info).RowsAffected == 0 {
		ctx.String(http.StatusBadRequest, "文件已失效")
	}
	info.DownloadNum++
	mysql.DB.Save(&info)

	pt := path.Join("static", "uploads", info.FileName)
	f, err := os.Open(pt)
	if err != nil {
		ctx.String(http.StatusBadRequest, "文件不存在")
		return
	}
	defer f.Close()

	p := make([]byte, 1024)

	w := ctx.Writer
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+url.QueryEscape(info.Name))

	var readErr error
	var readCount int

	for {
		readCount, readErr = f.Read(p)
		if readErr != nil {
			break
		}
		if readCount > 0 {
			if _, err := w.Write(p[:readCount]); err != nil {
				break
			}
		}
	}

	ctx.AbortWithStatus(http.StatusOK)
}

func Group(group *gin.RouterGroup) {
	group.GET("/download/:key", DownloadFile)
}
