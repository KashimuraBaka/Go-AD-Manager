package main

import "gitee.com/Kashimura/go-baka-control/cmd/web"

func main() {
	web.Router().Run(":80")
}
