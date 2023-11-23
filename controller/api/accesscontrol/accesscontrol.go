package accesscontrol

import (
	"net/http"
	"strconv"
	"time"

	"gitee.com/Kashimura/go-baka-control/db/mysql"
	g "gitee.com/Kashimura/go-baka-control/global"
	"github.com/gin-gonic/gin"
)

type ACUser struct {
	Userid      int       `json:"userid" gorm:"column:userid;primaryKey"`
	Badgenumber int       `json:"badgenid" gorm:"column:badgenumber"`
	Name        string    `json:"name" gorm:"column:name"`
	Group       int       `json:"group" gorm:"column:group"`
	Authority   int       `json:"authority" gorm:"column:authority"`
	Device      int       `json:"device" gorm:"column:device"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
	DeleteTime  time.Time `json:"delete_time" gorm:"column:delete_time"`
}

func (ACUser) TableName() string {
	return "ac_user"
}

type ACAuthority struct {
	Authority int    `json:"authority" gorm:"column:authority;primaryKey"`
	Name      string `json:"name" gorm:"column:name"`
}

func (ACAuthority) TableName() string {
	return "ac_authority"
}

type ACDevice struct {
	Device int    `json:"device" gorm:"column:device;primaryKey"`
	Name   string `json:"name" gorm:"column:name"`
	SN     string `json:"sn" gorm:"column:sn"`
}

func (ACDevice) TableName() string {
	return "ac_device"
}

type ACGroup struct {
	Group int    `json:"group" gorm:"column:group;primaryKey"`
	Name  string `json:"name" gorm:"column:name"`
}

func (ACGroup) TableName() string {
	return "ac_group"
}

func GetAuthorities(ctx *gin.Context) {
	authorities := make([]ACAuthority, 0)

	mysql.DB.Find(&authorities)

	ctx.JSON(http.StatusOK, g.ResponseSuccess(authorities, "success"))
}

func GetDevices(ctx *gin.Context) {
	devices := make([]ACDevice, 0)

	mysql.DB.Find(&devices)

	ctx.JSON(http.StatusOK, g.ResponseSuccess(devices, "success"))
}

func GetGroups(ctx *gin.Context) {
	groups := make([]ACGroup, 0)

	mysql.DB.Find(&groups)

	ctx.JSON(http.StatusOK, g.ResponseSuccess(groups, "success"))
}

func GetUsers(ctx *gin.Context) {
	users := make([]ACUser, 0)

	mysql.DB.Find(&users)

	ctx.JSON(http.StatusOK, g.ResponseSuccess(users, "success"))
}

func UserCreate(ctx *gin.Context) {
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
	user := &ACUser{}
	if deleteTime != "" {
		user.Badgenumber = badgenid
		user.Device = device
		user.Name = name
		user.Group = group
		user.Authority = authority
		user.CreateTime, _ = time.Parse("2006-01-02", createTime)
		user.DeleteTime, _ = time.Parse("2006-01-02", deleteTime)
		mysql.DB.Create(user)
		ctx.JSON(http.StatusOK, g.ResponseSuccess(user, "success"))
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
		ctx.JSON(http.StatusOK, g.ResponseSuccess(user, "susscess"))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseError("指纹已存在, 无法创建"))
	}
}

func UserUpdate(ctx *gin.Context) {
	userid, _ := strconv.Atoi(ctx.PostForm("userid"))
	if userid == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := &ACUser{Userid: userid}
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
		ctx.JSON(http.StatusOK, g.ResponseSuccess(user, "susscess"))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseError("指纹已存在, 无法创建"))
	}
}

func Group(group *gin.RouterGroup) {
	group.POST("/users", GetUsers)
	group.POST("/authorities", GetAuthorities)
	group.POST("/devices", GetDevices)
	group.POST("/groups", GetGroups)
	group.POST("/create", UserCreate)
	group.POST("/update", UserUpdate)
}
