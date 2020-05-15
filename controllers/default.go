package controllers

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Response(status int, data interface{}, err error) {

	if err != nil {
		c.Data["json"] = err.Error()

	} else {
		c.Data["json"] = data
	}
	c.Ctx.Output.SetStatus(status)
	c.ServeJSON()
	c.StopRun()
}
