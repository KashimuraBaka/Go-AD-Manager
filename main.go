package main

import "gitee.com/Kashimura/go-baka-control/cmd/gobaka"

func main() {
	gobaka.Router().Run(":80")
}
