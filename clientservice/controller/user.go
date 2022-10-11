package controller

import (
	"fmt"
	"net/http"
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/common/response"
	"seckill/dbservice/service"
	"sync"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	serv interfaces.UserServ
}

func (u *UserController) Get(ctx *gin.Context) {
	name := ctx.Param("name")
	user, err := u.serv.GetUser(name)
	if err != nil {
		fmt.Println(err.Error())
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, user, "ok")
}

func (u *UserController) Create(ctx *gin.Context) {
	user := &entity.User{}
	err := ctx.ShouldBindJSON(user)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	err = u.serv.AddUser(user)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	msg := fmt.Sprintf("create user %s success!", user.Name)
	response.Success(ctx, http.StatusOK, msg, "ok")
}

func (u *UserController) List(ctx *gin.Context) {
	users, err := u.serv.GetUsers()
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, users, "ok")

}

func (u *UserController) Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	err := u.serv.DeleteUser(name)
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", fmt.Sprintf("delete user %s success", name))

}

var userController *UserController

var userctlOnce = new(sync.Once)

func GetUserController() *UserController {
	userctlOnce.Do(func() {
		userController = &UserController{
			serv: service.GetUserServ(),
		}
	})
	return userController
}
