package middleware

import (
	"example/RestaurantProject/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			msg := fmt.Sprint("Something went wrong : No Authentication token provide")
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken) //genrate token claim struct instance using token string
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

		c.Next()

	}
}
