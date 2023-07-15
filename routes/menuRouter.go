package routes

import(
	"github.com/gin-gonic/gin"
	controller "example/RestaurantProject/controllers"
)

func MenuRoutes(incomingRoutes *gin.Engine){


	incomingRoutes.GET("/menus",controller.GetMenus())
	incomingRoutes.GET("/menus/:menu_id",controller.GetMenu())
	incomingRoutes.POST("/menus",controller.CreateMenus())
	incomingRoutes.PATCH("/menus/:menu_id",controller.UpdateMenu())

}