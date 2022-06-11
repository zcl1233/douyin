package controller

import (
	"net/http"
	"simple-demo-main/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	NP := NameAndPassword{}
	user := User{}
	var count int64 = 0
	models.DB.Model(&NP).Where("token=?", token).Count(&count)
	models.DB.Where("token=?", token).Find(&NP)
	models.DB.Where("name=?", NP.Name).Find(&user)
	if count == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
	} else if action_type == "1" {
		touserid, _ := strconv.ParseInt(to_user_id, 10, 64)
		relation := Relation{
			UserId:   user.Id,
			ToUserId: touserid,
		}
		models.DB.Create(&relation)
		touser := User{}
		models.DB.Where("id=?", touserid).Find(&touser)
		user.FollowCount += 1
		touser.FollowerCount += 1
		models.DB.Save(&user)
		models.DB.Save(&touser)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "follow successfully",
		})
	} else if action_type == "2" {
		touserid, _ := strconv.ParseInt(to_user_id, 10, 64)
		models.DB.Where("user_id=? AND to_user_id=?", user.Id, touserid).Delete(&Relation{})
		touser := User{}
		models.DB.Where("id=?", touserid).Find(&touser)
		user.FollowCount -= 1
		touser.FollowerCount -= 1
		models.DB.Save(&user)
		models.DB.Save(&touser)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "unfollow successfully",
		})
	}
}

// FollowList
func FollowList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("token=?", token).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
	} else {
		userid, _ := strconv.ParseInt(user_id, 10, 64)
		relations := []Relation{}
		models.DB.Where("user_id=?", userid).Find(&relations)
		touserid := make([]int64, len(relations))
		for _, relation := range relations {
			touserid = append(touserid, relation.ToUserId)
		}
		tousers := []User{}
		if len(touserid) != 0 {
			models.DB.Find(&tousers, touserid)
			for i := range tousers {
				tousers[i].IsFollow = true
			}
		}
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: tousers,
		})
	}
}

// FollowerList
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("token=?", token).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
	} else {
		userid, _ := strconv.ParseInt(user_id, 10, 64)
		relations := []Relation{}
		models.DB.Where("to_user_id=?", userid).Find(&relations)
		userList := []User{}
		users_id := make([]int64, len(relations))
		for _, relation := range relations {
			users_id = append(users_id, relation.UserId)
		}
		if len(users_id) != 0 {
			models.DB.Find(&userList, users_id)
			for i, user := range userList {
				var focus int64 = 0
				models.DB.Model(&Relation{}).Where("user_id=? AND to_user_id=?", userid, user.Id).Count(&focus)
				if focus == 1 {
					userList[i].IsFollow = true
				} else {
					userList[i].IsFollow = false
				}
			}
		}
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: userList,
		})
	}
}
