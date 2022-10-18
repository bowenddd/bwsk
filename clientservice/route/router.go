package route

import (
	"seckill/clientservice/controller"
	"seckill/clientservice/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Cors())
	r.Use(middleware.Auth())

	userGroup := r.Group("/user/")
	{

		userController := controller.GetUserController()

		userGroup.POST("login", userController.Login)

		userGroup.GET("name/:name", userController.Get)

		userGroup.POST("register", userController.Create)

		userGroup.GET("list", userController.List)

		userGroup.DELETE("name/:name", userController.Delete)
	}

	productGroup := r.Group("/product/")

	{
		productController := controller.GetProductController()

		productGroup.GET("name/:name", productController.Get)

		productGroup.GET("list", productController.List)

		productGroup.GET("stock/:id", productController.GetStock)

		productGroup.POST("create", productController.Create)

		productGroup.PUT("stock", productController.SetStock)

		productGroup.DELETE("/name:name", productController.Delete)

	}

	orderGroup := r.Group("/order/")

	{
		orderController := controller.GetOrderController()

		orderGroup.GET("id/:id", orderController.GetById)

		orderGroup.GET("uid/:uid", orderController.GetByUID)

		orderGroup.GET("pid/:pid", orderController.GetByPID)

		orderGroup.GET("list", orderController.List)

		orderGroup.POST("create", orderController.Create)

		orderGroup.DELETE("id/:id", orderController.Delete)

		orderGroup.PUT("clear", orderController.Clear)

	}

}
