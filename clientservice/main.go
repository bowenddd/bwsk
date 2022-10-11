package main

import (
	"github.com/gin-gonic/gin"
	"seckill/clientservice/route"
)

func main() {
	r := gin.Default()
	route.InitRouter(r)
	r.Run(":14643")
}
