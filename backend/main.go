package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/sirupsen/logrus"
	"imooc-product/backend/web/controllers"
	"imooc-product/common"
	"imooc-product/repositories"
	"imooc-product/services"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	//template := iris.HTML("./backend/web/views", ".html")
	template := iris.HTML("./backend/web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	//template := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	//4.设置模板目标
	app.HandleDir("/assets", "./backend/web/assets")

	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./backend/web/assets"))))
	//出现异常跳转到指定页面

	app.OnAnyErrorCode(func(ctx iris.Context) {
		fmt.Printf("code:%d", ctx.GetStatusCode())
		fmt.Println()
		logrus.Println("Before calling potentially problematic code...")
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		fmt.Println("调用shared/error.html")
		err := ctx.View("shared/error.html")
		if err != nil {
			logrus.Error(err)
		}
		logrus.Println("After calling potentially problematic code...")
	})

	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		logrus.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//5.注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productSerivce := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productSerivce)
	product.Handle(new(controllers.ProductController))
	//启动服务
	err = app.Run(iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
	if err != nil {
		return
	}

}
