package routes

import(
	"github.com/gin-gonic/gin"
	controller "example/RestaurantProject/controllers"
)

func OrderRoutes(incomingRoutes *gin.Engine){


	incomingRoutes.GET("/orders",controller.GetOrders())
	incomingRoutes.GET("/orders/:order_id",controller.GetOrder())
	incomingRoutes.POST("/orders",controller.CreateOrders())
	incomingRoutes.PATCH("/orders/:order_id",controller.UpdateOrder())

}