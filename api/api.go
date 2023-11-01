package api

import (
	"gitee.com/Kashimura/go-baka-control/api/dc"
	"gitee.com/Kashimura/go-baka-control/api/list"
	"gitee.com/Kashimura/go-baka-control/api/phone"
	"gitee.com/Kashimura/go-baka-control/api/user"
	"github.com/gin-gonic/gin"
)

func Api(group *gin.RouterGroup) {
	user.Group(group.Group("user"))
	list.Group(group.Group("list"))
	dc.Group(group.Group("dc"))
	phone.Group(group.Group("phone"))
}
