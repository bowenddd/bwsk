package route

import (
	"seckill/clientservice/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(Cors())
	userGroup := r.Group("/user/")

	{
		userController := controller.GetUserController()

		userGroup.GET(":name", userController.Get)

		userGroup.POST("create", userController.Create)

		userGroup.GET("list", userController.List)

		userGroup.DELETE(":name", userController.Delete)
	}

	productGroup := r.Group("/product/")

	{
		productController := controller.GetProductController()

		productGroup.GET(":name", productController.Get)

		productGroup.GET("list", productController.List)

		productGroup.GET("stock/:id", productController.GetStock)

		productGroup.POST("create", productController.Create)

		productGroup.PUT("stock", productController.SetStock)

		productGroup.DELETE(":name", productController.Delete)

	

	}

	orderGroup := r.Group("/order/")

	{
		orderController := controller.GetOrderController()

		orderGroup.GET("id/:id", orderController.GetById)

		orderGroup.GET("uid/:uid", orderController.GetByUID)

		orderGroup.GET("pid/:pid", orderController.GetByPID)

		orderGroup.GET("list", orderController.List)

		orderGroup.POST("create", orderController.Create)

		orderGroup.DELETE(":id", orderController.Delete)

		orderGroup.PUT("clear", orderController.Clear)

	}

}