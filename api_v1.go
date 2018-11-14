package main

import (
	"github.com/kataras/iris/mvc"
	"gotest/controllers"
)

func api_v1(app *mvc.Application) {

	// Create our movie repository with some (memory) data from the datasource.
	//repo := repositories.NewMovieRepository(datasource.Movies)
	// Create our movie service, we will bind it to the movie app's dependencies.
	//movieService := services.NewMovieService(repo)
	//app.Register(movieService)
	mvc.Configure(app.Party("/user"), user_v1)
}

func user_v1(app *mvc.Application){
	app.Handle(new(controllers.UserController))
}
