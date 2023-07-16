package controllers

import (
	"context"
	"example/RestaurantProject/database"
	"example/RestaurantProject/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {

			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage
		startIndex, err3 := strconv.Atoi(c.Query("startIndex"))
		if err3 != nil {
			log.Fatal(err3)
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: 1},
					{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			projectStage,
		})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error  occured while getting user list"})
			return
		}

		var allUsers []bson.M

		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		//return list of users
		c.JSON(http.StatusOK, allUsers)
	}

}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user) //find user matching to user_id and decode using struc user
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to found user"})
		}

		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		// convert the json data comming from postman to something that golang understands

		//validate the data based on user struct

		//you'll check if the email has already been used by another user

		//hash password

		//we will also check if the phone no. has already has been used by another user

		// create  some extra details - created_at , updated_at ,ID

		//genrate token and refersh token : (using func [genrateAllToken] from helper)

		// if all ok , then you insert this new user into the user collection

		//return status ok and send the result back

	}
}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		//unmarshal json data from postman

		// match user credential
		// if user exist

		//then we will veryfy password

		//if all goes well , then we will genrate tokens

		//update tokens - token and refresh token

		//return status okay

	}
}

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {

}

//we are using JWT authentication

/**     # SignUp  #

// convert the json data comming from postman to something that golang understands

	//validate the data based on user struct

	//you'll check if the email has already been used by another user

	//hash password

	//we will also check if the phone no. has already has been used by another user

	// create  some extra details - created_at , updated_at ,ID

	//genrate token and refersh token : (using func [genrateAllToken] from helper)

	// if all ok , then you insert this new user into the user collection

	//return status ok and send the result back
*/

/**    # LogIn  #

//unmarshal json data from postman

	// match user credential
	// if user exist

	//then we will veryfy password

	//if all goes well , then we will genrate tokens

	//update tokens - token and refresh token

	//return status okay
*/
