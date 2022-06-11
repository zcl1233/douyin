package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"simple-demo-main/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("token=?", token).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
	} else {
		title := c.PostForm("title")
		data, err := c.FormFile("data")
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}

		filename := filepath.Base(data.Filename)
		models.DB.Where("token=?", token).Find(&NP)
		user := User{}
		models.DB.Where("name=?", NP.Name).Find(&user)
		finalName := fmt.Sprintf("%d_%s", user.Id, filename)
		saveFile := filepath.Join("./public/", finalName)
		if err := c.SaveUploadedFile(data, saveFile); err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		videoname := strconv.FormatInt(user.Id, 10) + title + ".mp4"
		models.UploadAliyunOss(videoname, saveFile)
		video := Video{
			AuthorId:      user.Id,
			PlayUrl:       "https://ling-hu.oss-cn-shenzhen.aliyuncs.com/douyin/" + videoname,
			CoverUrl:      "",
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         title,
			PublishTime:   time.Now().Unix(),
		}
		models.DB.Model(&Video{}).Create(&video)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  finalName + " uploaded successfully",
		})
	}
}

// PublishList
func PublishList(c *gin.Context) {
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
		user := User{}
		models.DB.Where("token=?", token).Find(&NP)
		models.DB.Where("name=?", NP.Name).Find(&user)
		videoList := []Video{}
		models.DB.Where("author_id=?", user.Id).Find(&videoList)
		for i := range videoList {
			videoList[i].Author = user
		}
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: videoList,
		})
	}
}
