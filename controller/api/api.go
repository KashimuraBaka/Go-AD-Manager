package api

import (
	"gitee.com/Kashimura/go-baka-control/controller/api/dc"
	"gitee.com/Kashimura/go-baka-control/controller/api/list"
	"gitee.com/Kashimura/go-baka-control/controller/api/phone"
	"gitee.com/Kashimura/go-baka-control/controller/api/user"
	"gitee.com/Kashimura/go-baka-control/controller/file"
	"github.com/gin-gonic/gin"
)

func Group(group *gin.RouterGroup) {
	user.Group(group.Group("user"))
	list.Group(group.Group("list"))
	dc.Group(group.Group("dc"))
	phone.Group(group.Group("phone"))
	file.Group(group.Group("file"))
}
