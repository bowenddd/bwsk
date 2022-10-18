package controller

import (
	"fmt"
	"net/http"
	"seckill/authentication/jwt"
	"seckill/clientservice/service"
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/common/response"
	"seckill/common/vo"
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

func (u *UserController) Login(ctx *gin.Context) {
	userInfo := &vo.UserInfoVo{}
	err := ctx.ShouldBindJSON(userInfo)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user, err := u.serv.GetUser(userInfo.UserName)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if user.Password != userInfo.Password {
		response.Error(ctx, http.StatusOK, "password error")
		return
	}
	token, err := jwt.GenerateToken(userInfo.UserName,user.ID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, token, "ok")
}

var userController *UserController

var userctlOnce = new(sync.Once)

func GetUserController() *UserController {
	userctlOnce.Do(func() {
		userController = &UserController{
			serv: service.GetUserService(),
		}
	})
	return userController
}
