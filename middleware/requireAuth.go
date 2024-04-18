package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/timorodr/server/initializers"
	"github.com/timorodr/server/models"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("In Middleware")
	// Get the cookie off req
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// Decode/validate it

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// check the expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find the user with token sub
		// userID, ok := claims["sub"].(string)
		// if !ok {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		// 	return
		// }

		var user models.User
		client := initializers.ConnectToDb()
		collection := client.Database("jwtauth").Collection("users")
		collection.FindOne(context.Background(), bson.M{"_id": claims["sub"].(string)}).Decode(&user)

		// if user.ID ==  {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		// 	return
		//   }
		//   defer client.Disconnect(context.Background())
		// objectID, err := primitive.ObjectIDFromHex(userID)
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Error converting user ID"})
		// 	return
		// }
		// Attach to req

		// err = collection.FindOne(context.Background(), bson.M{"_id": claims["sub"].(string)}).Decode(&user)
		// if err != nil {
		// // 	// if err == mongo.ErrNoDocuments { // Handle user not found
		// // 	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// // 	// 	return
		// // 	// }
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
		// 	return
		// }

		c.Set("user", user)
		// Next should be used only inside middleware. It executes the pending handlers in the chain inside the calling handler
		// Continue
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
