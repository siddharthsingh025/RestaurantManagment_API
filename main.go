package main

import (
	"example/RestaurantProject/database"
	middleware "example/RestaurantProject/middleware"
	"example/RestaurantProject/routes"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food") /// this will creare database for me
//If a collection does not exist, MongoDB creates the collection when you first store data for that collection

func main() {

	// //set database url 
	// dburl := os. Getenv ("DATABASE_URL")
	// database.SetUrl(dburl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)
	routes.TableRoutes(router)
	router.Run(":" + port)

}
