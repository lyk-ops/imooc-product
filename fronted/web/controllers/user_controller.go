package controllers

import (
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"imooc-product/datamodels"
	"imooc-product/services"
	"imooc-product/tool"
	"strconv"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
	Session sessions.Sessions
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register",
	}
}

func (c *UserController) PostRegister() {
	var (
		nickName = c.Ctx.FormValue("nickName")
		password = c.Ctx.FormValue("password")
		userName = c.Ctx.FormValue("userName")
	)
	user := &datamodels.User{UserName: userName, NickName: nickName, HashPassword: password}
	_, err := c.Service.AddUser(user)
	if err != nil {
		c.Ctx.Redirect("/user/error")
		return
	}
	c.Ctx.Redirect("/user/login")

}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login",
	}
}

func (c *UserController) PostLogin() mvc.Response {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)
	user, ok := c.Service.IsPwdSuccess(username, password)
	if !ok {
		return mvc.Response{
			Path: "/user/error",
		}
	}
	tool.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	return mvc.Response{
		Path: "/product",
	}
}
