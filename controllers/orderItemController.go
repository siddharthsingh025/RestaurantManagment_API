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

type OrderItemPack struct {
	Table_id   *string
	Order_item []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderItemCollection.Find(context.TODO(), bson.M{}) //bson.M{}  is a empty query on monodb
		defer cancel()                                                    //context. TODO when it's unclear which Context to use or it is not yet available
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing orderItems"})
			return
		}

		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allOrderItems)

	}

}

func GetOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderItemId := c.Param("order_item_id")

		var orderItem models.OrderItem

		err := orderCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		//here we call findOne function that will return only one record  with matching orderItemId and after getting  we cast it into "&orderItem" - orderItemModel

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching orderItem"})
			return
		}
		//return back to user with status code 200
		c.JSON(http.StatusOK, orderItem)

	}
}

// this function will return all orderItems belongs to that particular order which is assosiated order_id we are passing in request
func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId) // take the orderId and give all the orderItems in that order
		// a order consist of multiple orderItems

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items by order ID"})
			return
		}

		//return back to user with status code 200
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func CreateOrderItems() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var orderItemPack OrderItemPack
		var order models.Order

		//binding/cast recieved  json data with order model
		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))

		orderItemToBeInserted := []interface{}{}
		order.Table_id = orderItemPack.Table_id  // here we actually accessing table id of orderItem via orderItemPack
		order_id := OrderItemOrderCreator(order) // generate new order and return its order_id

		for _, orderItem := range orderItemPack.Order_item { //range over orderItems associated with that particular order_id

			orderItem.Order_id = order_id

			//to validate data that we get from user
			//it uses tag we provide in our models to validate data { validator pkg}
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			//create orderItem  of orderId - > "order_id"
			orderItem.ID = primitive.NewObjectID()
			//timeStamps
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()

			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num

			orderItemToBeInserted = append(orderItemToBeInserted, orderItem) // add into list of orderItems
		}

		result, orderItemInsertErr := orderItemCollection.InsertMany(ctx, orderItemToBeInserted)

		if orderItemInsertErr != nil {
			msg := fmt.Sprintf("orderItems were not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		//return back to user with status code 200
		c.JSON(http.StatusOK, result)

	}
}

func UpdateOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var orderItem models.OrderItem

		orderItemId := c.Param("order_item_id")

		//binding/cast recieved  json data with order model
		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"order_item_id": orderItemId}
		var updateObj primitive.D

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{Key: "unit_price", Value: *orderItem.Unit_price})

		}

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: *orderItem.Quantity})

		}

		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{Key: "food_id", Value: *orderItem.Food_id})

		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint("orderItem updation failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		//return back to user with status code 200
		c.JSON(http.StatusOK, result)

	}
}

//primitive
// M is an unordered representation of a BSON document.
//This type should be used when the order of the elements does not matter

// order -> has multiple -> orderItems and -> every OrderItem has food -> Food

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}                                                                               // gives all the orderItems of that particular order_id
	lookupFoodStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}} // lookup on food data of matched OrderItems

	//$lookup is tool to used for looking up other collection in database via foreignKey we have in our collection
	//let say orderItem want to lookup data of Food using food_id store in both ;
	// localField : in orderItem data  {food_id}
	// ForeignField : in Food data     {food_id}

	unwindFoodStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}
	//unwindStage used to unwind array that we got from lookupStage so that we can  access & perform more actions on it
	//path is the field which we want to unwind
	//preserveNullAndEmptyArrays : to preserve all the fields that are empty or null in array

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"locakField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	//used to controll which data whill go ahead and which are not
	projectStage := bson.D{

		{
			"$project", bson.D{
				{"id", 0},
				{"amount", "$food.price"},
				{"total_count", 1},
				{"food_name", "$food.name"},
				{"food_image", "$food.food_image"},
				{"table_number", "$table.table_number"},
				{"table_id", "$table.table_id"},
				{"order_id", "$order.order_id"},
				{"price", "$food.price"},
				{"quantity", 1},
			}}}

	//grouping data
	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"order_items", bson.D{{"$push", "$$ROOT"}}}}}}

	projectStage2 := bson.D{
		{
			"$project", bson.D{
				{"id", 0},
				{"payment_due", 1},
				{"table_count", 1},
				{"table_number", "$_id.table_number"},
				{"order_items", 1},
			}}}

	// order is mendatory and very important as we know its a pipe line so output from first stage is used in next and things are hase be to done in same order
	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupFoodStage,
		unwindFoodStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})

	if err != nil {
		panic(err)
	}

	if err := result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	defer cancel()

	return OrderItems, err

}
