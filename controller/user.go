package controller

import (
	"net/http"
	"simple-demo-main/models"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	newuser := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&newuser).Where("name=?", username).Count(&count)
	if count == 1 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User already exist",
		})
	} else {
		token := username + password
		newNameAndPassword := NameAndPassword{
			Name:     username,
			Password: password,
			Token:    token,
		}
		newUser := User{
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		models.DB.Create(&newNameAndPassword)
		models.DB.Create(&newUser)
		models.DB.Where("name=?", username).Find(&LoginUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   LoginUser.Id,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("name=? AND password=?", username, password).Count(&count)
	if count != 1 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User does not exist",
		})
	} else {
		user := User{}
		models.DB.Where("name=?", username).Find(&user)
		LoginUser = user
		token := username + password
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	user := User{}
	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("token=?", token).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
		})
	} else {
		models.DB.Where("token=?", token).Find(&NP)
		models.DB.Where("name=?", NP.Name).Find(&user)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	}
}
