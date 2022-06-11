package controller

import (
	"net/http"
	"simple-demo-main/models"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed
func Feed(c *gin.Context) {
	token := c.Query("token")
	//latest_time:=c.Query("latest_time")
	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("Token=?", token).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
	} else {
		videos := []Video{}
		models.DB.Limit(10).Order("publish_time").Find(&videos)
		for i, video := range videos {
			authorId := video.AuthorId
			author := User{}
			models.DB.Where("id=?", authorId).Find(&author)
			models.DB.Where("token=?", token).Find(&NP)
			loger := User{}
			models.DB.Where("name=?", NP.Name).Find(&loger)
			var focus int64 = 0
			models.DB.Model(&Relation{}).Where("user_id=? AND to_user_id=?", loger.Id, authorId)
			if focus == 1 {
				author.IsFollow = true
			} else {
				author.IsFollow = false
			}
			videos[i].Author = author
			var like int64 = 0
			models.DB.Model(&Favorite{}).Where("video_id=? AND user_name=?", video.Id, loger.Name).Count(&like)
			if like == 1 {
				videos[i].IsFavorite = true
			} else {
				videos[i].IsFavorite = false
			}
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videos,
			NextTime:  time.Now().Unix(),
		})
	}
}
