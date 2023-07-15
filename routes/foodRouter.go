package routes

import (
	controller "example/RestaurantProject/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/foods", controller.GetFoods())
	incomingRoutes.GET("/food/:food_id", controller.GetFood())
	incomingRoutes.POST("/foods", controller.CreateFoods())
	incomingRoutes.PATCH("/food/:food_id", controller.UpdateFood())

}
