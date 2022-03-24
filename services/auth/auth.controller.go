package auth

import (
	"context"
	"fmt"
	"net/http"
	JWTManager "swiggy/gin/lib/helpers"
	db "swiggy/gin/lib/utils"
	"swiggy/gin/services/user"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type loginBody struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type signUpBody struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Login(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body := loginBody{}

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	user := &user.User{}
	if err := db.DataStore.Collection("user").FindOne(ctx, bson.M{"username": body.UserName}).Decode(&user); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User Not found"})
	}

	if user == nil || !user.IsCorrectPassword(body.Password) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Password Incorrect"})
	}

	token, err := JWTManager.Manager.Generate(user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error while creating token"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"user": user, token: token})
}

func Signup(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body := signUpBody{}

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	user, err := user.NewUser(body.UserName, body.Password, body.Role)

	res, err := db.DataStore.Collection("user").InsertOne(ctx, user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("can not convert to oid %v", err)})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"username": body.UserName, "Id": oid.Hex()})
}

func CheckAuth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Auth Working"})
}
