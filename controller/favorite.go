package controller

import (
	"net/http"
	"simple-demo-main/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("Token=?", token).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
	} else {
		video := Video{}
		var videoid, _ = strconv.ParseInt(video_id, 10, 64)
		models.DB.Where("id=?", videoid).Find(&video)
		models.DB.Where("Token=?", token).Find(&NP)
		if action_type == "1" {
			video.FavoriteCount += 1
			video.IsFavorite = true
			models.DB.Save(&video)
			favorite := Favorite{
				VideoId:  videoid,
				UserName: NP.Name,
			}
			models.DB.Create(&favorite)
		} else if action_type == "2" {
			video.FavoriteCount -= 1
			if video.FavoriteCount == 0 {
				video.IsFavorite = false
			}
			models.DB.Save(&video)
			models.DB.Where("user_name=?", NP.Name).Delete(&Favorite{})
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
		})
	}
}

// FavoriteList
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("Token=?", token).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
	} else {
		models.DB.Where("Token=?", token).Find(&NP)
		favoriteList := []Favorite{}
		models.DB.Where("user_name=?", NP.Name).Find(&favoriteList)
		video_id := make([]int64, len(favoriteList))
		for _, favorite := range favoriteList {
			video_id = append(video_id, favorite.VideoId)
		}
		videoList := []Video{}
		if len(video_id) != 0 {
			models.DB.Find(&videoList, video_id)
			for i, video := range videoList {
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
				videoList[i].Author = author
				videoList[i].IsFavorite = true
			}
		}
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: videoList,
		})
	}
}
