package response

type Sign struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     uint64 `json:"user_id"`
	Token      string `json:"token"`
}
