package response

type User struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint64 `json:"follow_count"`
	FollowerCount uint64 `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type Info struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	User       *User  `json:"user"`
}
