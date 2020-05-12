package main

import (
	_ "media_service/routers"
	services "media_service/services/db"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	services.InitSql()
	beego.Run()
}
