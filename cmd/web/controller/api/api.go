package api

import (
	"gitee.com/Kashimura/go-baka-control/cmd/web/controller/api/accesscontrol"
	"gitee.com/Kashimura/go-baka-control/cmd/web/controller/api/attendance"
	"gitee.com/Kashimura/go-baka-control/cmd/web/controller/api/dc"
	"gitee.com/Kashimura/go-baka-control/cmd/web/controller/api/user"
	"github.com/gin-gonic/gin"
)

func Register(parentGroup *gin.Engine) {
	group := parentGroup.Group("/api")

	// 注册用户服务
	user.Register(group)

	accesscontrol.Register(group)
	attendance.Register(group)
	dc.Register(group)
}
