package controllers

import (
	"context"
	"example/RestaurantProject/database"
	"example/RestaurantProject/models"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodb "food" collection instance
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

var validate = validator.New()

//Validate is designed to be thread-safe and used as a singleton instance. It caches information about your struct and validations, in essence only parsing your validation tags once per struct type. Using multiple instances neglects the benefit of caching

func GetFoods() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage")) // to limit the sending data to frontEnd max recordPerPage set to be 10
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10 // default limit
		}

		page, err := strconv.Atoi(c.Query("page")) // getting  query parameter keyed with "page"
		if err != nil || page < 1 {
			page = 1 // default limit
		}

		startIndex := (page - 1) * recordPerPage // where to start from
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		//aggregation pipeline stage and operator : mongodb-topic

		// will fo to databse and try to match with data
		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}} // match using notthing

		// will group all the records by group by name , in group stage on counting data our data get lost for that we create $data
		groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}}, {Key: "total_count", Value: bson.D{{Key: "$sum ,1"}}}, {Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}

		// Reshapes a document stream by renaming, adding, or removing fields.
		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},         // "0" means this field won't go to frontEnd
					{Key: "total_count", Value: 1}, // "1" means this field can  go to frontEnd
					{Key: "food_item", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}, // now we get above data here as &data
					// $slice will controlls number of elements to return
					//startIndex : where to start from
					// recordPerPage : how many pages do you want to sent
				},
			},
		}

		// oerder of stages  is  important here
		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			groupStage,
			projectStage,
		})

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occuered while listing food item"})
		}

		//
		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allFoods[0])

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
		//return back to user with status code 200
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

		//timeStamps
		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		//inserting into database
		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintf("food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return

		}
		defer cancel()
		//return back to user with status code 200
		c.JSON(http.StatusOK, result)
	}
}

func UpdateFood() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food

		foodId := c.Param("food_id")

		//binding/cast recieved  json data with food model
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price})
		}

		if food.Food_image != nil {
			updateObj = append(updateObj, bson.E{Key: "food_image", Value: food.Food_image})

		}

		if food.Menu_id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
			defer cancel()
			if err != nil {
				msg := fmt.Sprint("message: Menu was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.Menu_id}) // doubt
		}

		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.Updated_at})

		upsert := true

		filter := bson.M{"food_id": foodId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint("food item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

func round(num float64) int {

	return int(num + math.Copysign(0.5, num))

}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

/**
Query returns the keyed url query value if it exists, otherwise it returns an empty string `("")`. It is shortcut for `c.Request.URL.Query().Get(key)`

    GET /path?id=1234&name=Manu&value=
       c.Query("id") == "1234"
       c.Query("name") == "Manu"
       c.Query("value") == ""
       c.Query("wtf") == ""
*/
