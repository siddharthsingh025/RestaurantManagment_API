package controllers

import (
	"context"
	"example/RestaurantProject/database"
	"example/RestaurantProject/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMenus() gin.HandlerFunc {

	var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*(time.Second))

		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the menu items"})
		}

		var allMenu []bson.M
		if err = result.All(ctx, &allMenu); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenu)
	}

}

func GetMenu() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		//here we call findOne function that will find menu with id "menuId" and after getting it we cast it into "&menu" - MenuModel
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the menu from databse"})
		}

		c.JSON(http.StatusOK, menu)

	}
}

func CreateMenus() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		//binding/cast recieved  json data with food model
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//to validate data that we get from user
		//it uses tag we provide in our models to validate data { validator pkg}
		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		//timeStamps
		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		//inserting into database
		result, insertErr := menuCollection.InsertOne(ctx, menu)

		if insertErr != nil {
			msg := fmt.Sprint("menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		//return back to user with status code 200
		c.JSON(http.StatusOK, result)
		defer cancel()

	}
}

func UpdateMenu() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		//binding/cast recieved  json data with menu model
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//getting id from request context
		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj primitive.D

		if menu.Start_Date != nil && menu.End_Date != nil {

			//to validate is give start time is after current time and end time if after start
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
				msg := "kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}
			// appending single entitiiy to updated map
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})

			if menu.Name != "" {
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
			}

			if menu.Category != "" {
				updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
			}

			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

			// MongoDB, upsert is an option that is used for update operation e.g. update(), findAndModify(), etc.
			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}
			// updating data into database
			result, err := menuCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{"$set", updateObj},
				},
				&opt,
			)

			if err != nil {
				msg := "Menu update failed"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			}

			//return back to user with status code 200
			defer cancel()
			c.JSON(http.StatusOK, result)
		}
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(check) && end.After(start)
}

/**

BSON stands for Binary Javascript Object Notation.
It is a binary-encoded serialization of JSON documents.
BSON has been extended to add some optional non-JSON-native data types, like dates and binary data.
BSON can be compared to other binary formats, like Protocol Buffers.
*/

/**

MongoDB stores documents in a binary representation called BSON that allows for easy and flexible data processing.

The Go Driver provides 4 main types for working with BSON data:

// D: An ordered representation of a BSON document (slice)

// M: An unordered representation of a BSON document (map)

// A: An ordered representation of a BSON array

// E: A single element inside a D type

The following example demonstrates how to construct a query filter using the bson.D type to match documents with a quantity field value greater than 100:

filter := bson.D{{"quantity", bson.D{{"$gt", 100}}}}
*/
