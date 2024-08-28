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

// NoCacheMiddleware 是一个中间件，用于设置缓存控制头
func NoCacheMiddleware(ctx iris.Context) {
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")
	ctx.Next()
}
func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("info")
	//3.注册模板

	//template := iris.HTML("./backend/web/views", ".html")
	template := iris.HTML("./backend/web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	//template := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template.Reload(true))
	//4.设置模板目标
	app.HandleDir("/assets", "./backend/web/assets")

	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./backend/web/assets"))))
	//出现异常跳转到指定页面

	app.OnAnyErrorCode(func(ctx iris.Context) {
		fmt.Printf("code:%d\n", ctx.GetStatusCode())
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		if err := ctx.View("shared/error.html"); err != nil {
			app.Logger().Errorf("Failed to render error page: %v", err)
		}
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
	productParty := app.Party("/product", NoCacheMiddleware)
	product := mvc.New(productParty)
	product.Register(ctx, productSerivce)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order", NoCacheMiddleware)
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))
	//启动服务
	err = app.Run(iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
	if err != nil {
		return
	}

}
