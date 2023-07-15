package controllers

import (
	"context"
	"example/RestaurantProject/database"
	"example/RestaurantProject/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongodb "food" collection instance
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

var validate = validator.New()

//Validate is designed to be thread-safe and used as a singleton instance. It caches information about your struct and validations, in essence only parsing your validation tags once per struct type. Using multiple instances neglects the benefit of caching

func GetFoods() gin.HandlerFunc {

	return func(c *gin.Context) {

	}

}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		//here we call findOne function that will find food with id "foodId" and after getting it we cast it into "&food" - Foodmodel
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching food from databse"})
		}

		c.JSON(http.StatusOK, food)

	}
}

func CreateFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food

		//binding/cast recieved  json data with food model
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//to validate data that we get from user
		//it uses tag we provide in our models to validate data { validator pkg}
		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		//checking where the user sent food is present in our menu or not
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()

		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		food.Created_at, _ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintf("food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return

		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func UpdateFood() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}

func round(num float64) int {

}

func toFixed(num float64, precision int) float64 {

}
