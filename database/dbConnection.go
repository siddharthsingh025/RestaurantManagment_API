package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodb://YOUR_USERNAME_HERE:YOUR_PASSWORD_HERE@0.0.0.0:YOUR_LOCALHOST_PORT_HERE/
// my mongodb  cluster url :mongodb+srv://siddharthsingh25102001:5Lc6HWebOZ5Of6a8@cluster0.hd8ampj.mongodb.net/?retryWrites=true&w=majority

// var url string = ""
// func SetUrl(link string) {
// 	url = link 
// }

func DBinstance() *mongo.Client {

	MongoDb_url := os. Getenv ("DATABASE_URL")
		if MongoDb_url ==""{
			log.Fatal("databse url invalid")
			MongoDb_url = "mongodb+srv://siddharthsingh25102001:5Lc6HWebOZ5Of6a8@cluster0.hd8ampj.mongodb.net/?retryWrites=true&w=majority"
		}
	fmt.Print(MongoDb_url)

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb_url))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("connected to mongodb")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)

	return collection
}

/**
By passing a context.Context across function calls and services in your application, those can stop working early and return an error when their processing is no longer needed. For more about Context, see Go Concurrency Patterns: Context.

For example, you might want to:

End long-running operations, including database operations that are taking too long to complete.
Propagate cancellation requests from elsewhere, such as when a client closes a connection.


You can use a Context to set a timeout or deadline after which an operation will be canceled. To derive a Context with a timeout or deadline, call context.WithTimeout or context.WithDeadline.
  // Create a Context with a timeout.
    queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
*/

/*
A collection is a grouping of MongoDB documents.
 Documents within a collection can have different fields.
 A collection is the equivalent of a table in a relational database system.
  A collection exists within a single database.
*/
