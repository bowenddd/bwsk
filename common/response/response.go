package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Success(ctx *gin.Context, code int, data interface{}, msg string) {
	rps := Response{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	ctx.JSON(code, rps)
}

func Error(ctx *gin.Context, code int, msg string) {
	rps := Response{
		Code: code,
		Msg:  msg,
	}
	fmt.Println(rps)
	ctx.JSON(code, rps)
}
