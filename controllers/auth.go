package controllers

import (
	"go-backend/models"
	"go-backend/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// DeleteUser godoc
// @Summary      Deletes The user with id
// @Description  Responds with the deletion of user
// @Tags         Delete
// @Produce      json
// @Router       /user/delete/:id [POST]
func DeleteUser(c *gin.Context) {
	var user []models.User
	id := c.Param("id")
	if err := models.DB.Where("id=?", id).First(&user).Delete(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted successfully",
	})
}

// AllUser       godoc
// @Summary      Return all users
// @Description  Responds with the data of users
// @Tags         AllUser
// @Produce      json
// @Router       /users [GET]
func AllUser(c *gin.Context) {
	var user []models.User
	models.DB.Find(&user)
	c.JSON(http.StatusOK, &user)
}

// GetUser godoc
// @Summary      Get the user with id
// @Description  Responds with the user details 
// @Tags         GetUser
// @Produce      json
// @Router       /user/:id [GET]
func GetUser(c *gin.Context) {
	var user []models.User
	if err := models.DB.Where("id=?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "record not found",
		})
	}
	c.JSON(http.StatusOK, &user)
}

// Register      godoc
// @Summary      Register The user
// @Description  Responds with the Registeration of user
// @Tags         Register
// @Produce      json
// @Router       /register [POST]
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	u := models.User{}
	u.Username = input.Username
	u.Password = input.Password
	_, err := u.SaveUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "registered successfully",
	})
}

// CurrentUser godoc
// @Summary      Shows the current user logged in.
// @Description  respond with json of user data
// @Tags         Current User
// @Produce      json
// @Router       /admin/user [GET]
func CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login         godoc
// @Summary      Login view
// @Description  Login the user
// @Tags         Login
// @Produce      json
// @Router       /login [POST]
func Login(context *gin.Context) {
	var input LoginInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserByUsername(input.Username)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidatePass(input.Password)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := token.GenrateToken(int(user.ID))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	context.JSON(http.StatusOK, gin.H{"token": jwt})
}
