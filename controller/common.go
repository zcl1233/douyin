package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	AuthorId      int64  `json:"author_id"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
	PublishTime   int64  `json:"publish_time,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	VideoId    int64  `json:"video_id"`
	UserId     int64  `json:"user_id"`
	User       User   `json:"users"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type NameAndPassword struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Favorite struct {
	VideoId  int64  `json:"video_id"`
	UserName string `json:"user_name"`
}
type Relation struct {
	UserId   int64 `json:"user_id"`
	ToUserId int64 `json:"to_user_id"`
	User     User  `json:"user,omitempty"`
}

func (User) TableName() string {
	return "users"
}
func (NameAndPassword) TableName() string {
	return "nameandpasswords"
}
func (Video) TableName() string {
	return "videos"
}
func (Favorite) TableName() string {
	return "favorites"
}
func (Comment) TableName() string {
	return "comments"
}
func (Relation) TableName() string {
	return "relations"
}
