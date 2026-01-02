package user

import (
	"net/http"
	"strconv"
	"time"

	"gitee.com/Kashimura/go-baka-control/services/db/mysql"
	"gitee.com/Kashimura/go-baka-control/utils/buffer"
	"gitee.com/Kashimura/go-baka-control/utils/webhttp"
	"gitee.com/Kashimura/go-baka-control/utils/webhttp/jwt"
	"github.com/gin-gonic/gin"
)

func Register(parentGroup *gin.RouterGroup) {
	group := parentGroup.Group("/user")

	group.POST("/loginService", loginService)
	group.POST("/refresh", refreshToken)
	group.GET("/record", GetUserRecord)
	group.GET("/record_page", GetUserRecordPage)
}

func loginService(ctx *gin.Context) {
	params := LoginParmas{}
	ctx.ShouldBindJSON(&params)

	// 密码或账号为空则返回非法请求
	if params.Username == "" || params.Password == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// 查询是否有该账户
	user := &mysql.UserLogon{UserName: params.Username}
	// 如果账户存在且密码异样
	if mysql.DB.Find(user).RowsAffected != 0 && string(buffer.Base64Decode(params.Password)) == user.PassWord {
		data := make(map[string]any)
		data["token"], _ = jwt.GenUserToken(0, user.UserName, 2*time.Hour)
		data["lastvisit"] = user.LogonTime
		data["lastip"] = user.IP
		// 更新数据
		user.IP = ctx.RemoteIP()
		user.LogonTime = time.Now()
		mysql.DB.Save(user)
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Data: data,
			Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
		})
	} else {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_WRONG_PASSWORD,
			Msg:  webhttp.GetMessageByCode(webhttp.ERROR_WRONG_PASSWORD),
		})
	}
}

func refreshToken(ctx *gin.Context) {
	data := UserLogonInfo{}
	data.IP = ctx.RemoteIP()

	// 解析Token
	authorization := ctx.Request.Header.Get("Authorization")
	uc, err := jwt.ParseUserToken(authorization)
	if err != nil {
		record := mysql.UserRecord{IP: data.IP}
		// 记录用户访问网页时间
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
		data.ReadTime = time.Now()
		data.ReadNum = record.ReadNum
		// 返回结果
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_AUTH_CHECK_TOKEN_FAIL,
			Data: data,
			Msg:  webhttp.GetMessageByCode(webhttp.ERROR_AUTH_CHECK_TOKEN_FAIL),
		})
	} else {
		user := mysql.UserLogon{UserName: uc.UserName}
		// 如果用户不存在则直接返回错误
		if mysql.DB.Find(&user).RowsAffected == 0 {
			ctx.JSON(http.StatusOK, webhttp.Response{
				Code: webhttp.ERROR_UNKNOWN_USER,
				Msg:  webhttp.GetMessageByCode(webhttp.ERROR_UNKNOWN_USER),
			})
			return
		}
		// 更新记录
		user.LogonTime = time.Now()
		user.IP = ctx.RemoteIP()
		mysql.DB.Save(&user)
		// 返回结果
		data.Name = user.UserName
		data.IP = user.IP
		data.ReadTime = user.LogonTime
		data.Token, err = jwt.RefreshUserToken(authorization)
		if err != nil {
			ctx.JSON(http.StatusOK, webhttp.Response{
				Code: webhttp.ERROR_OTHER,
				Msg:  "Token 出现异常: " + err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, webhttp.Response{
				Code: webhttp.SUCCESS,
				Data: data,
				Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
			})
		}
	}
}

func GetUserRecord(ctx *gin.Context) {
	users := &[]mysql.UserRecord{}
	if mysql.DB.Find(users).RowsAffected == 0 {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "not found",
		})
	} else {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Data: users,
			Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
		})
	}
}

func GetUserRecordPage(ctx *gin.Context) {
	var total int64
	var users []mysql.UserRecord
	pageNum, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("size"))
	tx := mysql.DB.Model(&mysql.UserRecord{})
	tx.Count(&total)
	if tx.Offset((pageNum-1)*pageSize).Limit(pageSize).Find(&users).RowsAffected == 0 {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.ERROR_OTHER,
			Msg:  "not found",
		})
	} else {
		ctx.JSON(http.StatusOK, webhttp.Response{
			Code: webhttp.SUCCESS,
			Data: gin.H{"total": total, "data": users},
			Msg:  webhttp.GetMessageByCode(webhttp.SUCCESS),
		})
	}
}
