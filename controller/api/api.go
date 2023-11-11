package api

import (
	"gitee.com/Kashimura/go-baka-control/controller/api/accesscontrol"
	"gitee.com/Kashimura/go-baka-control/controller/api/dc"
	"gitee.com/Kashimura/go-baka-control/controller/api/list"
	"gitee.com/Kashimura/go-baka-control/controller/api/phone"
	"gitee.com/Kashimura/go-baka-control/controller/api/user"
	"gitee.com/Kashimura/go-baka-control/controller/file"
	"github.com/gin-gonic/gin"
)

func Group(group *gin.RouterGroup) {
	list.Group(group.Group("list"))
	file.Group(group.Group("file"))
	phone.Group(group.Group("phone"))
	user.Group(group.Group("user"))
	dc.Group(group.Group("dc"))
	accesscontrol.Group(group.Group("accesscontrol"))
}
