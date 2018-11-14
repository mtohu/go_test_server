package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"gotest/services"
	"strconv"
)

type UserController struct {
    Ctx iris.Context
    userServer services.IUserService
}

func (c *UserController) PostLogin(){

    c.Ctx.HTML("233333")
}

func (c *UserController) GetUserInfo(){
	c.userServer = services.NewUserService()
	//c.Ctx.Values().GetString()
	b := c.Ctx.Request().FormValue("b")
	a,err := strconv.Atoi(b)
	if (err !=nil){
		fmt.Println(a)
	}
	c.userServer.GetByUsernameAndPassword(b,b)
	//c.Ctx.WriteString("===============sss")
	c.Ctx.HTML("===="+b)
}

func (c *UserController) PostUserInfo(){

}

func (c *UserController) PutUserInfo(){

}

func (c *UserController) DeleteUser(){

}
