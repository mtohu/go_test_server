package routes

import  (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"gotest/controllers"
)

type UserRouter struct {
	 uparty iris.Party
}
func (u *UserRouter) SetUserRouter(app *iris.Application, path string) {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Authorization"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})
	u.uparty= app.Party(path,crs).AllowMethods(iris.MethodOptions)
	//路由分发,这里再次路由分发，将功能块再次细化
	u.setLoginRoute()
	u.setUserInfoRoute()
}
/*
* 登录模块
* @uri:/mcos/login
*/
func (u *UserRouter) setLoginRoute() {
	// POST: http://localhost:8080/api/v1/login/
	u.uparty.Post("/login", func(ctx iris.Context) {

		hander_req_post := &controllers.UserController{
			Ctx: ctx,
		}
		hander_req_post.PostLogin()
	})
}

/*
* 用户信息处理模块路由
* 也是功能模块的入口(请求的控制器、服务处理和数据模型封装不在此说明)
* @uri:/mcos/userinfo
*/
func (u *UserRouter) setUserInfoRoute() {
	// GET: http://localhost:8080/api/v1/userinfo/42
	u.uparty.Get("/userinfo/{id:string}", func(ctx iris.Context) {
		hander_req_get := &controllers.UserController{
			Ctx: ctx,
		}
		hander_req_get.GetUserInfo()
	})
	// POST: http://localhost:8080/api/v1/userinfo/
	u.uparty.Post("/userinfo", func(ctx iris.Context) {
		hander_req_post := &controllers.UserController{
			Ctx: ctx,
		}
		hander_req_post.PostUserInfo()
	})
	// PUT: http://localhost:8080/api/v1/userinfo/
	u.uparty.Put("/userinfo/{id:int}", func(ctx iris.Context) {
		hander_req_put := &controllers.UserController{
			Ctx: ctx,
		}
		hander_req_put.PutUserInfo()
	})
	// DELETE: http://localhost:8080/api/v1/userinfo/42
	u.uparty.Delete("/userinfo/{id:int}", func(ctx iris.Context) {
		hander_req_del := &controllers.UserController{
			Ctx: ctx,
		}
		hander_req_del.DeleteUser()
	})
}

