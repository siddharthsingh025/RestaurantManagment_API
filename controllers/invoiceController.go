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

type InvoiceViewFormat struct {
	Invoice_id       string
	Order_id         string
	Payment_method   string
	Payment_status   *string
	Payment_due_date time.Time
	Payment_due      interface{}
	Table_number     interface{}
	Order_details    interface{}
}

// mongodb "order" collection instance
var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := invoiceCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing invoices"})
		}

		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allInvoices)

	}

}

func GetInvoice() gin.HandlerFunc {

	return func(c *gin.Context) {
		//creating context with request timeOut limit
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//getting id from request context
		invoiceId := c.Param("invoice_id")

		var invoice models.Invoice

		// finding weather invoice is present in record or not
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice) // case it into invoiceModel
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing invoices"})
		}

		//now intialise invoiceView parameter from invoiceModel object
		var invoiceView InvoiceViewFormat
		allOrderItems, err := ItemsByOrder(invoice.Order_id) // getting list or data off all order present in invoice
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date
		invoiceView.Payment_method = "null"
		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}

		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"] //now intialise invoiceView parameter from list of orders that we fetched above
		invoiceView.Table_number = allOrderItems[0]["total_number"]
		invoiceView.Order_details = allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK, invoiceView) // sent back invoiceView to user with ok(200) status

	}
}

func CreateInvoices() gin.HandlerFunc {

	return func(c *gin.Context) {

		//creating context with request timeOut limit
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var invoice models.Invoice

		//binding/cast recieved  json data with order model
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
		defer cancel()

		if err != nil {
			msg := fmt.Sprintf("order was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}

		invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))

		//timeStamps
		invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		//to validate data that we get from user
		//it uses tag we provide in our models to validate data { validator pkg}
		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		//inserting into database
		result, insertErr := invoiceCollection.InsertOne(ctx, invoice)

		if insertErr != nil {
			msg := fmt.Sprintf("invoice was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		//return back to user with status code 200
		c.JSON(http.StatusOK, result)

	}
}

func UpdateInvoice() gin.HandlerFunc {

	return func(c *gin.Context) {

		//creating context with request timeOut limit
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var invoice models.Invoice
		var invoiceId = c.Param("invoice_id")

		//binding/cast recieved  json data with order model
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if invoice.Payment_method != nil {

			updateObj = append(updateObj, bson.E{Key: "payment_method", Value: invoice.Payment_method})

		}

		if invoice.Payment_status != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_status", Value: invoice.Payment_status})
		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.Updated_at})

		filter := bson.M{"invoice_id": invoiceId}
		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint("invoice update  get failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		//return back to user with status code 200
		c.JSON(http.StatusOK, result)

	}
}

/**
M is an unordered representation of a BSON document.
This type should be used when the order of the elements does not matter.
This type is handled as a regular map[string]interface{} when encoding and decoding.
Elements will be serialized in an undefined, random order. If the order of the elements matters, a D should be used instead.

*/
