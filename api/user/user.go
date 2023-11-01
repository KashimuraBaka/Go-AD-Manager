package user

import (
	"net/http"
	"strconv"
	"time"

	"gitee.com/Kashimura/go-baka-control/db/mdb"
	"gitee.com/Kashimura/go-baka-control/db/mysql"
	g "gitee.com/Kashimura/go-baka-control/global"
	"gitee.com/Kashimura/go-baka-control/services/jwt"
	"gitee.com/Kashimura/go-baka-control/services/tools"
	"github.com/gin-gonic/gin"
)

func UserLogin(ctx *gin.Context) {
	username := ctx.PostForm("user")
	password := ctx.PostForm("pwd")
	// 密码或账号为空则返回非法请求
	if username == "" || password == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// 查询是否有该账户
	user := &UserLogon{UserName: username}
	if mysql.DB.Find(user).RowsAffected != 0 && string(tools.Base64Decode(password)) == user.PassWord {
		data := make(map[string]any)
		// 生成 Token
		data["token"] = jwt.GetToken(&jwt.PayLoad{SUB: username, IP: ctx.RemoteIP()}, 120)
		data["lastvisit"] = user.LogonTime
		data["lastip"] = user.IP
		// 更新数据
		user.IP = ctx.RemoteIP()
		user.LogonTime = time.Now()
		mysql.DB.Save(user)
		ctx.JSON(http.StatusOK, g.Response{Code: g.RES_SUCCESS, Data: data, Msg: "登陆成功"})
	} else {
		ctx.JSON(http.StatusOK, g.ResponseError("密码错误"))
	}
}

func UserCheck(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	payload, err := jwt.VerifyToken(authorization)
	data := make(map[string]any)
	data["ip"] = ctx.RemoteIP()
	if err != nil {
		record := UserRecord{IP: ctx.RemoteIP()}
		if mysql.DB.Find(&record).RowsAffected == 0 {
			record.ReadNum = 1
			record.ReadTime = time.Now()
			record.Role = 0
			mysql.DB.Create(&record)
		} else {
			record.ReadNum += 1
			record.ReadTime = time.Now()
			mysql.DB.Save(&record)
		}
		data["time"] = time.Now()
		data["num"] = record.ReadNum
		ctx.JSON(http.StatusOK, g.Response{Code: g.RES_ERROR, Data: data, Msg: "invalid token"})
	} else {
		sub := payload.SUB
		user := UserLogon{UserName: sub}
		if mysql.DB.Find(&user).RowsAffected == 0 {
			ctx.JSON(http.StatusOK, g.ResponseError("未找到该账户!"))
			return
		}
		data["user"] = user.UserName
		data["lastvisit"] = user.LogonTime
		data["lastip"] = user.IP
		user.LogonTime = time.Now()
		user.IP = ctx.RemoteIP()
		mysql.DB.Save(&user)
		ctx.JSON(http.StatusOK, g.ResponseSuccess(data))
	}
}

func GetUserRecord(ctx *gin.Context) {
	users := &[]UserRecord{}
	if mysql.DB.Find(users).RowsAffected == 0 {
		ctx.JSON(http.StatusOK, g.ResponseError("not found"))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseSuccess(users))
	}
}

func GetUserRecordPage(ctx *gin.Context) {
	var total int64
	var users []UserRecord
	pageNum, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("size"))
	tx := mysql.DB.Model(&UserRecord{})
	tx.Count(&total)
	if tx.Offset((pageNum-1)*pageSize).Limit(pageSize).Find(&users).RowsAffected == 0 {
		ctx.JSON(http.StatusOK, g.ResponseError("not found"))
	} else {
		ctx.JSON(http.StatusOK, g.ResponseSuccess(gin.H{
			"total": total,
			"data":  users,
		}))
	}
}

func InsertAttendanceRecord(ctx *gin.Context) {
	user := ctx.PostForm("user")
	time := ctx.PostForm("time")
	device := ctx.PostForm("device")
	// 空数据不处理
	if user == "" || time == "" || device == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// 数值转为整数
	userid, err := strconv.Atoi(user)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	// 查询是否有该用户
	if mdb.DB.CheckUserID(userid) != nil {
		ctx.JSON(http.StatusOK, g.ResponseError("未找到该用户"))
		return
	}

	if err := mdb.DB.InsertCheckInOut(userid, time, device); err != nil {
		ctx.JSON(http.StatusBadRequest, g.ResponseError("数据插入失败"))
		return
	}

	ctx.JSON(http.StatusOK, g.ResponseSuccess(nil))
}

func Group(group *gin.RouterGroup) {
	group.POST("/login", UserLogin)
	group.POST("/check", UserCheck)
	group.GET("/record", GetUserRecord)
	group.GET("/record_page", GetUserRecordPage)
	group.POST("/attendance", InsertAttendanceRecord)
}
