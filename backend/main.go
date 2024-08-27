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
	"log"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	tmplate := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板目标
	app.HandleDir("/assets", "./backend/web/assets", iris.DirOptions{
		// 这里可以添加一些选项，但如果你不需要特别配置，可以省略这些选项
		// 例如，启用或禁用浏览静态目录（ListDirectory: true/false）
		// 自定义HTTP头（Headers: map[string]string{...}）
		// 启用或禁用索引文件（IndexName: "index.html"）等
	})

	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./backend/web/assets"))))
	//出现异常跳转到指定页面

	app.OnAnyErrorCode(func(ctx iris.Context) {
		log.Println("Before calling potentially problematic code...")
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		err := ctx.View("shared/error.html")
		if err != nil {
			fmt.Println(err)
			log.Println("After calling potentially problematic code...")
			return
		}
		log.Println("After calling potentially problematic code...")
	})
	// 第 32 行及其周围的代码

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
	app.Run(iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
