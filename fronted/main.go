package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/sirupsen/logrus"
	"imooc-product/common"
	"imooc-product/fronted/web/controllers"
	"imooc-product/repositories"
	"imooc-product/services"
	"time"
)

func main() {
	// 1. 创建 iris 实例
	app := iris.New()
	// 2. 设置错误模式，在 mvc 模式下提示错误
	app.Logger().SetLevel("debug")

	// 注册模板
	template := iris.HTML("./fronted/web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(template)

	// 4. 设置模板目标
	app.HandleDir("/public", "./fronted/web/public")

	// 出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		fmt.Printf("code:%d\n", ctx.GetStatusCode())
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		if err := ctx.View("user/register"); err != nil {
			app.Logger().Errorf("Failed to render error page: %v", err)
		}
	})

	// 连接数据库
	fmt.Printf("connecting to database...\n")
	db, err := common.NewMysqlConn()
	if err != nil {
		logrus.Error(err)
		fmt.Printf("connect database failed\n")
		return // 遇到错误后退出程序
	}
	defer db.Close()

	// 注册控制器
	sess := sessions.New(sessions.Config{
		Cookie:  "AdminCookie",
		Expires: 600 * time.Minute,
	})

	user := repositories.NewUserRepository("user", db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))

	// 这里修正了传递的参数
	userPro.Register(userService, context.Background(), sess.Start)

	userPro.Handle(new(controllers.UserController))

	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
