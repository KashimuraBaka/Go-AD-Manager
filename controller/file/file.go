package file

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"gitee.com/Kashimura/go-baka-control/db/mysql"
	g "gitee.com/Kashimura/go-baka-control/global"
	"gitee.com/Kashimura/go-baka-control/services/jwt"
	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context) {
	// 获取上传文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, g.ResponseError("上传失败, 非有效文件!"))
		return
	}
	suffix := path.Ext(file.Filename)
	// 获取文件MD5值
	reader, _ := file.Open()
	t, _ := io.ReadAll(reader)
	fileMD5 := fmt.Sprintf("%x%s", md5.Sum(t), suffix)
	// 查询数据库是否存在该MD5值
	rfile := FileInfo{}
	if mysql.DB.Model(rfile).Where("fname = ?", fileMD5).Find(&rfile).RowsAffected == 0 {
		// 保存文件
		ctx.SaveUploadedFile(file, path.Join("static", "uploads", fileMD5))
		// 保存记录
		rfile.Name = file.Filename
		rfile.FileName = fileMD5
		rfile.Size = file.Size // DefaultSize: b
		rfile.User = ctx.RemoteIP()
		rfile.Date = time.Now()
		mysql.DB.Create(&rfile)
		// filename 改为 JWT密钥 提供下载
		rfile.FileName = jwt.GetFileToken(rfile.ID)
		ctx.JSON(http.StatusOK, g.ResponseSuccess(rfile, "上传成功"))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseError("文件已存在!"))
	}
}

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

func ReNameFile(ctx *gin.Context) {
	key, name := ctx.PostForm("key"), ctx.PostForm("name")
	// 如果 key 为空的话
	if key == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// 校验 文件token
	id, err := jwt.VerifyFileToken(key)
	if err != nil {
		ctx.JSON(http.StatusOK, g.ResponseError("非有效文件密钥"))
		return
	}
	// 查询该文件
	file := FileInfo{ID: id}
	if mysql.DB.Find(&file).RowsAffected == 0 {
		ctx.JSON(http.StatusOK, g.ResponseError("重命名失败, 请刷新网页重试!"))
		return
	}
	file.ReName = name
	mysql.DB.Save(&file)
	ctx.JSON(http.StatusOK, g.ResponseSuccess(nil, "重命名成功!"))
}

func DeleteFile(ctx *gin.Context) {
	key := ctx.PostForm("key")
	// 如果 key 为空的话
	if key == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// 校验文件token
	id, err := jwt.VerifyFileToken(key)
	if err != nil {
		ctx.JSON(http.StatusOK, g.ResponseError("非有效文件密钥"))
		return
	}
	// 数据库查找文件
	file := FileInfo{ID: id}
	if mysql.DB.Find(&file).RowsAffected == 0 {
		ctx.JSON(http.StatusOK, g.ResponseError("资源文件不存在"))
		return
	}
	// 删除文件
	filepath := path.Join("static", "uploads", file.FileName)
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		if err = os.Remove(filepath); err != nil {
			ctx.JSON(http.StatusOK, g.ResponseError("文件删除失败, 请稍后重试"))
			return
		}
	}
	mysql.DB.Delete(&file)
	ctx.JSON(http.StatusOK, g.ResponseSuccess(nil, "删除成功"))
}

func DeleteBeforeFile(ctx *gin.Context) {
	type tFileInfo struct {
		ID          int64     `json:"id" gorm:"column:id"`
		Name        string    `json:"name" gorm:"column:name"`
		FileName    string    `json:"fname" gorm:"column:fname"`
		ReName      string    `json:"rname" gorm:"column:rname"`
		Size        int64     `json:"size" gorm:"column:size"`
		User        string    `json:"user" gorm:"column:user"`
		Date        time.Time `json:"time" gorm:"column:time"`
		DownloadNum int       `json:"download_num" gorm:"column:downloadnum"`
		Role        string    `json:"role" gorm:"column:role"`
	}

	// 查找历史文件
	files := make([]tFileInfo, 0)
	res := mysql.DB.Table(
		"list_file file, user_record user, user_role role",
	).Where(
		"file.user = user.ip and user.role = role.role_id and role.role_name != 'admin' and NOW() - INTERVAL 7 DAY > file.time",
	).Select(
		"id, name, fname, rname, size, user, time, downloadnum, IFNULL(role_name, 'user') AS role",
	).Find(&files)

	// 开始删除文件
	if res.RowsAffected > 0 {
		var row map[string]any
		data := make([]map[string]any, 0)
		for _, v := range files {
			filepath := path.Join("static", "uploads", v.FileName)
			if _, err := os.Stat(filepath); !os.IsNotExist(err) {
				if err = os.Remove(filepath); err != nil {
					continue
				} else {
					row = make(map[string]any)
					row["id"] = v.ID
					row["name"] = v.Name
					data = append(data, row)
				}
			}
		}
		ctx.JSON(http.StatusOK, g.ResponseSuccess(data, "删除历史文件完毕"))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseSuccess(nil, "无过期文件需要删除"))
	}
}

func Group(group *gin.RouterGroup) {
	group.POST("/upload", UploadFile)
	group.POST("/rename", ReNameFile)
	group.POST("/delete", DeleteFile)
	group.POST("/deletebefore", DeleteBeforeFile)
	group.GET("/download/:key", DownloadFile)
}
