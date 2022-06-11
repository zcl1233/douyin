package controller

import (
	"simple-demo-main/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"net/http"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}
type CommentResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("token=?", token).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
	} else if action_type == "1" {
		comment_text := c.Query("comment_text")
		var videoid, _ = strconv.ParseInt(video_id, 10, 64)
		models.DB.Where("Token=?", token).Find(&NP)
		user := User{}
		models.DB.Where("name=?", NP.Name).Find(&user)
		comment := Comment{
			VideoId:    videoid,
			UserId:     user.Id,
			Content:    comment_text,
			CreateDate: time.Now().Format("01-02"),
		}
		models.DB.Create(&comment)
		video := Video{}
		models.DB.Where("id", videoid).Find(&video)
		video.CommentCount += 1
		models.DB.Save(&video)
		comment.User = user
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 0,
			},
			Comment: comment,
		})
	} else if action_type == "2" {
		comment_id := c.Query("comment_id")
		commentid, _ := strconv.ParseInt(comment_id, 10, 64)
		comment := Comment{}
		models.DB.Where("id=?", commentid).Find(&comment)
		models.DB.Where("id=?", commentid).Delete(&Comment{})
		video := Video{}
		var videoid, _ = strconv.ParseInt(video_id, 10, 64)
		models.DB.Where("id", videoid).Find(&video)
		video.CommentCount -= 1
		models.DB.Save(&video)
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 0,
			},
			Comment: comment,
		})
	}
}

// CommentList
func CommentList(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	NP := NameAndPassword{}
	var count int64 = 0
	models.DB.Model(&NP).Where("token=?", token).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user doesn't exist",
		})
	} else {
		var videoid, _ = strconv.ParseInt(video_id, 10, 64)
		comments := []Comment{}
		models.DB.Where("video_id=?", videoid).Find(&comments)
		for index, comment := range comments {
			comment_userid := comment.UserId
			comment_user := User{}
			models.DB.Where("id=?", comment_userid).Find(&comment_user)
			video := Video{}
			models.DB.Where("id=?", videoid).Find(&video)
			loger := User{}
			models.DB.Where("token=?", token).Find(&NP)
			models.DB.Where("name=?", NP.Name).Find(&loger)
			var focus int64 = 0
			models.DB.Model(&Relation{}).Where("user_id=? AND to_user_id=?", loger.Id, comment_user.Id).Count(&focus)
			if focus == 1 {
				comment_user.IsFollow = true
			} else {
				comment_user.IsFollow = false
			}
			comments[index].User = comment_user
		}
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{
				StatusCode: 0,
			},
			CommentList: comments,
		})
	}
}
