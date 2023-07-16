package controllers

import (
	"context"
	"example/RestaurantProject/database"
	helpers "example/RestaurantProject/helpers"
	"example/RestaurantProject/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		// convert the json data comming from postman to something that golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//validate the data based on user struct
		//to validate data that we get from user
		//it uses tag we provide in our models to validate data { validator pkg}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		//you'll check if the email has already been used by another user

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		//hash password

		password := HashPassword(*user.Password)
		user.Password = &password

		//we will also check if the phone no. has already has been used by another user
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exsists"})
			return
		}

		// create  some extra details - created_at , updated_at ,ID

		//timeStamps
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		//genrate token and refersh token : (using func [genrateAllToken] from helper)

		token, refreshToken, _ := helpers.GenrateAllToken(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		// if all ok , then you insert this new user into the user collection

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("user were not created: something went wrong")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		//return status ok and send the result back
		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		//unmarshal json data from postman
		// convert the json data comming from postman to something that golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// match user credential

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found : recheck your credential "})
			return
		}

		//then we will veryfy password

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		//if all goes well , then we will genrate tokens
		token, refreshToken, _ := helpers.GenrateAllToken(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		//update tokens - token and refresh token
		helpers.UpdateAllToken(token, refreshToken, foundUser.User_id)

		//return status okay
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(providePassword), []byte(userPassword))

	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprint("user or passward is incorrect")
		check = false
	}
	return check, msg
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
