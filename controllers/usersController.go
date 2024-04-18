package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // make sure v5 is used at endpoint
	"github.com/timorodr/server/initializers"
	"github.com/timorodr/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)


func HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}


func Signup(c *gin.Context) {
	//Get the email and password off req body
	var body struct {
		Email    string
		Password string
	}

	// get the above variables off request body retrun to stop
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	// hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	// create the user
	user := models.User{Email: body.Email, Password: string(hash)}
	client := initializers.ConnectToDb()

	collection := client.Database("jwtauth").Collection("users")

	// Insert the user document into the collection
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		// Handle potential errors like duplicate email
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already exists",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func Login(c *gin.Context) {
	// Get the email and password off req body
	var body struct {
		Email    string
		Password string
	}

	//checks the Method and Content-Type to select a binding engine automatically, Depending on the "Content-Type" header different bindings are used
	// application/json
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	client := initializers.ConnectToDb()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to DB"})
	// 	return
	// }
	// defer client.Disconnect(context.Background())
	// look up requested user
	collection := client.Database("jwtauth").Collection("users")

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": body.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
		return
	}
	// compare credentials with saved user credentials
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}
	//Gen JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}
	// send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true) // look into changing arguments here false -> true for deployment?

	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context){
	user, _ := c.Get("user")



	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}