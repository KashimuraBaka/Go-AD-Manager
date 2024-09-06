package accesscontrol

import (
	"net/http"
	"strconv"
	"time"

	"gitee.com/Kashimura/go-baka-control/cmd/web/controller/middleware"
	"gitee.com/Kashimura/go-baka-control/services/db/mysql"
	"gitee.com/Kashimura/go-baka-control/utils/webhttp"
	"github.com/gin-gonic/gin"
)

func Register(parentGroup *gin.RouterGroup) {
	group := parentGroup.Group("/ac")

	group.Use(middleware.VerifyAccountToken())
	group.POST("/users", getUsers)
	group.POST("/authorities", getAuthorities)
	group.POST("/devices", getDevices)
	group.POST("/groups", getGroups)
	group.POST("/create", createUser)
	group.POST("/update", updateUser)
}

func getAuthorities(ctx *gin.Context) {
	authorities := make([]mysql.ACAuthority, 0)

	mysql.DB.Find(&authorities)

	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Data: authorities,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}

func getDevices(ctx *gin.Context) {
	devices := make([]mysql.ACDevice, 0)

	mysql.DB.Find(&devices)

	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}

func getGroups(ctx *gin.Context) {
	groups := make([]mysql.ACGroup, 0)

	mysql.DB.Find(&groups)

	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}

func getUsers(ctx *gin.Context) {
	users := make([]mysql.ACUser, 0)

	mysql.DB.Find(&users)

	ctx.JSON(http.StatusOK, webhttp.Response{
		Code: webhttp.SUCCESS,
		Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
	})
}

func createUser(ctx *gin.Context) {
	badgenid, _ := strconv.Atoi(ctx.PostForm("badgenid"))
	name := ctx.PostForm("name")
	group, _ := strconv.Atoi(ctx.PostForm("group"))
	authority, _ := strconv.Atoi(ctx.PostForm("authority"))
	device, _ := strconv.Atoi(ctx.PostForm("device"))
	createTime := ctx.PostForm("create_time")
	deleteTime := ctx.PostForm("delete_time")
	if badgenid == 0 || name == "" || group == 0 || authority == 0 || device == 0 || createTime == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := &mysql.ACUser{}
	if deleteTime != "" {
		user.Badgenumber = badgenid
		user.Device = device
		user.Name = name
		user.Group = group
		user.Authority = authority
		user.CreateTime, _ = time.Parse("2006-01-02", createTime)
		user.DeleteTime, _ = time.Parse("2006-01-02", deleteTime)
		mysql.DB.Create(user)
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
		})
	} else if mysql.DB.Model(user).Where(
		"(delete_time IS NULL OR delete_time = ?) AND badgenumber = ? And device = ?", "0000-00-00", badgenid, device,
	).Find(user).RowsAffected == 0 {
		user.Badgenumber = badgenid
		user.Device = device
		user.Name = name
		user.Group = group
		user.Authority = authority
		user.CreateTime, _ = time.Parse("2006-01-02", createTime)
		mysql.DB.Create(user)
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
		})
	} else {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "指纹已存在, 无法创建",
		})
	}
}

func updateUser(ctx *gin.Context) {
	userid, _ := strconv.Atoi(ctx.PostForm("userid"))
	if userid == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := &mysql.ACUser{Userid: userid}
	if mysql.DB.Find(user).RowsAffected != 0 {
		// 名字
		if name := ctx.PostForm("name"); name != "" {
			user.Name = name
		}
		// 指纹
		if badgenid, _ := strconv.Atoi(ctx.PostForm("badgenid")); badgenid > 0 {
			user.Badgenumber = badgenid
		}
		// 组
		if group, _ := strconv.Atoi(ctx.PostForm("group")); group > 0 {
			user.Group = group
		}
		// 权限
		if authority, _ := strconv.Atoi(ctx.PostForm("authority")); authority > 0 {
			user.Authority = authority
		}
		// 设备
		if device, _ := strconv.Atoi(ctx.PostForm("device")); device > 0 {
			user.Device = device
		}
		// 创建时间
		if createTime := ctx.PostForm("create_time"); createTime != "" {
			user.CreateTime, _ = time.Parse("2006-01-02", createTime)
		}
		// 删除时间
		if deleteTime := ctx.PostForm("delete_time"); deleteTime != "" {
			user.DeleteTime, _ = time.Parse("2006-01-02", deleteTime)
		}
		mysql.DB.Save(user)
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
		})
	} else {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "指纹已存在, 无法创建",
		})
	}
}
