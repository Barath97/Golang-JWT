package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jwt/database"
	"github.com/jwt/helpers"
	"github.com/jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// HashPassword hashes the password using bcrypt
func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

// VerifyPassword checks if the provided password matches the hashed password
func VerifyPassword(hashedPassword, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, "password is incorrect"
	}
	return true, ""
}

// Signup handles user registration
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the email or phone already exists
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil || count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil || count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "phone number already exists"})
			return
		}

		// Hash the password and create the user
		hashedPassword := HashPassword(*user.Password)
		user.Password = &hashedPassword
		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = primitive.NewObjectID()

		// Convert user ID to a string and set User_id
		userId := user.ID.Hex()
		user.User_id = &userId

		// Generate JWT tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, userId)
		user.Token = &token
		user.Refresh_token = &refreshToken

		// Insert user into database
		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": user})
	}
}

// Login handles user authentication
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find user by email
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// Verify password
		passwordIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		// Generate new JWT tokens and update user
		token, refreshToken, _ := helpers.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.First_name,
			*foundUser.Last_name,
			*foundUser.User_type,
			*foundUser.User_id, // Dereference the pointer here
		)
		helpers.UpdateAllTokens(token, refreshToken, *foundUser.User_id) // Dereference here as well

		c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token, "refresh_token": refreshToken})
	}
}

// GetUsers retrieves all users (admin-only)
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("user_type") != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var users []models.User
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving users"})
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var user models.User
			if err := cursor.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding user"})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

// GetUser retrieves a single user by ID
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var user models.User

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
