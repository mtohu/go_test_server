package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"golang.org/x/net/netutil"
	"gotest/datasource"
	"gotest/routes"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"time"
)

func main() {
	_,err := datasource.Getinstance()
	if(err !=nil){
		log.Println("init database  failure...")
		os.Exit(1)
	}

	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Logger().SetLevel("debug")
	app.RegisterView(iris.HTML("./web/views", ".html"))
	app.OnErrorCode(iris.StatusNotFound, notFounds)
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// close all hosts
		app.Shutdown(ctx)
	})
	// Serve our controllers.
	new(routes.UserRouter).SetUserRouter(app, "/api/v1")
	l, err := net.Listen("tcp4", ":8188")
	if err != nil {
		panic(err)
	}
	l = netutil.LimitListener(l, 1000)//设置开启协层数量
	app.Run(iris.Listener(l),iris.WithoutInterruptHandler,iris.WithConfiguration(iris.YAML("./configs/app.yml")))
	//go app.Run(iris.Addr(":8188"),iris.WithoutInterruptHandler,iris.WithConfiguration(iris.YAML("./configs/app.yml")))
	//app.NewHost(&http.Server{Addr:":8187"}).ListenAndServe()
}
func notFounds(ctx iris.Context) {
	// 当http.status=400 时向客户端渲染模板$views_dir/errors/404.html
	/*req := HttpRequest.NewRequest()
	res, err := req.Get("http://127.0.0.1:8000?id=10&title=HttpRequest",nil)
	body, err := res.Body()
	if err != nil {
		log.Println(err)
		return
	}*/
	//fmt.Println(ctx.)
	ctx.View("errors/404.html")
}
func noConnectDatabases(ctx iris.Context) {
	// 当http.status=400 时向客户端渲染模板$views_dir/errors/404.html
	/*req := HttpRequest.NewRequest()
	res, err := req.Get("http://127.0.0.1:8000?id=10&title=HttpRequest",nil)
	body, err := res.Body()
	if err != nil {
		log.Println(err)
		return
	}*/
	//fmt.Println(ctx.)
	ctx.View("errors/404.html")
}



