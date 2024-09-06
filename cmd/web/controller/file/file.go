package file

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"gitee.com/Kashimura/go-baka-control/services/db/mysql"
	"gitee.com/Kashimura/go-baka-control/utils/webhttp"
	"gitee.com/Kashimura/go-baka-control/utils/webhttp/jwt"
	"github.com/gin-gonic/gin"
)

func Register(parentGroup *gin.Engine) {
	group := parentGroup.Group("/file")

	group.GET("/list", getFiles)
	group.POST("/upload", uploadFile)
	group.POST("/rename", renameFile)
	group.POST("/delete", deleteFile)
	group.POST("/deletebefore", deleteBeforeFile)
	group.GET("/download/:key", downloadFile)
}

func getFiles(ctx *gin.Context) {
	res := []mysql.FileInfo{}

	if mysql.DB.Order("time DESC").Find(&res).RowsAffected > 0 {
		for i, v := range res {
			res[i].FileName, _ = jwt.GenFileToken(v.ID, 6*time.Hour)
		}
	}

	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}

func uploadFile(ctx *gin.Context) {
	// 获取上传文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "上传失败, 非有效文件!",
		})
		return
	}
	suffix := path.Ext(file.Filename)
	// 获取文件MD5值
	reader, _ := file.Open()
	t, _ := io.ReadAll(reader)
	fileMD5 := fmt.Sprintf("%x%s", md5.Sum(t), suffix)
	// 查询数据库是否存在该MD5值
	rfile := mysql.FileInfo{}
	if mysql.DB.Model(rfile).Where("fname = ?", fileMD5).Find(&rfile).RowsAffected == 0 {
		// 保存文件
		ctx.SaveUploadedFile(file, path.Join("data", "uploads", fileMD5))
		// 保存记录
		rfile.Name = file.Filename
		rfile.FileName = fileMD5
		rfile.Size = file.Size // DefaultSize: b
		rfile.User = ctx.RemoteIP()
		rfile.Date = time.Now()
		mysql.DB.Create(&rfile)
		// filename 改为 JWT密钥 提供下载
		rfile.FileName, _ = jwt.GenFileToken(rfile.ID, 6*time.Hour)
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Msg:  "上传成功",
		})
	} else {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "文件已存在!",
		})
	}
}

func downloadFile(ctx *gin.Context) {
	// 获取文件ID
	fc, err := jwt.ParseFileToken(ctx.Param("key"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// 查询记录
	info := mysql.FileInfo{ID: fc.FileID}
	if mysql.DB.Find(&info).RowsAffected == 0 {
		ctx.String(http.StatusBadRequest, "文件已失效")
		return
	}
	info.DownloadNum++
	mysql.DB.Save(&info)

	// 打开文件
	filePath := path.Join("data", "uploads", info.FileName)

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		ctx.String(http.StatusNotFound, "File not found")
		return
	}

	// 获取文件大小
	fileSize := fileInfo.Size()

	// 打开文件并将其内容写入响应体中
	file, err := os.Open(filePath)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer file.Close()

	// 获取文件MIME类型
	buffer := make([]byte, 512)
	_, _ = file.Read(buffer)
	file.Seek(0, io.SeekStart)

	// 设置响应头
	ctx.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(info.Name))
	ctx.Header("Content-Type", http.DetectContentType(buffer))
	ctx.Header("Accept-Ranges", "bytes")

	// 获取请求头中的 Range 字段
	rangeHeader := ctx.Request.Header.Get("Range")

	if rangeHeader == "" {
		ctx.Header("Content-Length", strconv.FormatInt(fileSize, 10))
		io.Copy(ctx.Writer, file)
	} else {
		var start, end int64
		if _, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		if start < 0 || end < 0 || start >= fileSize || end >= fileSize || start > end {
			ctx.String(http.StatusBadRequest, "Invalid Range header")
			return
		}

		ctx.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		ctx.Header("Content-Length", strconv.FormatInt(start-end, 10))
		ctx.Status(http.StatusPartialContent)

		// 读取位置
		_, err = file.Seek(start, io.SeekStart)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// 将文件内容写入响应体中
		io.CopyN(ctx.Writer, file, end-start)
	}
}

func renameFile(ctx *gin.Context) {
	key, name := ctx.PostForm("key"), ctx.PostForm("name")
	// 如果 key 为空的话
	if key == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 校验 文件token
	fc, err := jwt.ParseFileToken(key)
	if err != nil {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "非有效文件密钥",
		})
		return
	}

	// 查询该文件
	file := mysql.FileInfo{ID: fc.FileID}
	if mysql.DB.Find(&file).RowsAffected == 0 {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "重命名失败, 请刷新网页重试!",
		})
		return
	}

	file.ReName = name
	mysql.DB.Save(&file)
	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Msg:  "重命名成功!",
	})
}

func deleteFile(ctx *gin.Context) {
	key := ctx.PostForm("key")
	// 如果 key 为空的话
	if key == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 校验文件token
	fc, err := jwt.ParseFileToken(key)
	if err != nil {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "非有效文件密钥",
		})
		return
	}

	// 数据库查找文件
	file := mysql.FileInfo{ID: fc.FileID}
	if mysql.DB.Find(&file).RowsAffected == 0 {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "资源文件不存在",
		})
		return
	}

	// 删除文件
	filepath := path.Join("data", "uploads", file.FileName)
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		if err = os.Remove(filepath); err != nil {
			ctx.JSON(http.StatusOK, webhttp.Response{
				Code: webhttp.ERROR_OTHER,
				Msg:  "文件删除失败, 请稍后重试",
			})
			return
		}
	}
	mysql.DB.Delete(&file)
	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Msg:  "删除成功",
	})
}

func deleteBeforeFile(ctx *gin.Context) {
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
			filepath := path.Join("data", "uploads", v.FileName)
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
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Msg:  "删除历史文件完毕",
		})
	} else {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Msg:  "无过期文件需要删除",
		})
	}
}
