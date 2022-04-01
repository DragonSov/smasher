package handlers

import (
	"fmt"
	"github.com/DragonSov/smasher/server/domain/Users"
	"github.com/DragonSov/smasher/server/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandlers struct {
	Service services.UserService
}

func (h UserHandlers) CreateUser(c *gin.Context) {
	// Creating a user model
	u := Users.UserModel{}

	// Checking form fields
	u.Login, u.Password = c.PostForm("login"), c.PostForm("password")
	if u.Login == "" || u.Login == " " || u.Password == "" || u.Password == " " {
		c.JSON(406, gin.H{
			"status":  "error",
			"message": "The form fields are not filled in",
		})
		return
	}

	// Selecting this user
	user, err := h.Service.SelectUserByLogin(u.Login)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "An unexpected error has occurred",
		})
		return
	} else if user != nil {
		c.JSON(409, gin.H{
			"status":  "error",
			"message": "The user already exists",
		})
		return
	}

	// Creating new user
	user, err = h.Service.CreateUser(u)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "An unexpected error has occurred",
		})
		return
	}

	// Return user login
	c.JSON(200, gin.H{
		"status": "success",
		"login":  user.Login,
	})
}

func (h UserHandlers) SelectUserByUUID(c *gin.Context) {
	// Parsing user ID
	userID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to get a user ID. Try again later",
		})
		return
	}

	// Selecting this user
	user, err := h.Service.SelectUserByUUID(userID)
	if err != nil || user == nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	// Return information about this user
	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func (h UserHandlers) SelectUserByLogin(c *gin.Context) {
	// Selecting user by login
	user, err := h.Service.SelectUserByLogin(c.Query("login"))
	if err != nil || user == nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	// Return information about this user
	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func (h UserHandlers) SignIn(c *gin.Context) {
	// // Checking form fields
	login, password := c.PostForm("login"), c.PostForm("password")
	if login == "" || login == " " || password == "" || password == " " {
		c.JSON(406, gin.H{
			"status":  "error",
			"message": "The form fields are not filled in",
		})
		return
	}

	// User authentication
	token, err := h.Service.SignIn(login, password)
	if err != nil {
		c.JSON(403, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Return access token
	c.JSON(200, gin.H{
		"status":       "success",
		"access_token": token,
	})
}
