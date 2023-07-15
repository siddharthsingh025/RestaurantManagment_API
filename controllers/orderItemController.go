package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderItems() gin.HandlerFunc {

	return func(ctx *gin.Context) {

	}

}

func GetOrderItem() gin.HandlerFunc {

	return func(ctx *gin.Context) {

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {

}

func CreateOrderItems() gin.HandlerFunc {

	return func(ctx *gin.Context) {

	}
}

func UpdateOrderItem() gin.HandlerFunc {

	return func(ctx *gin.Context) {

	}
}

//primitive
// M is an unordered representation of a BSON document.
//This type should be used when the order of the elements does not matter
