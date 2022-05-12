package response

type SignUp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserID     uint   `json:"user_id"`
	Token      string `json:"token"`
}
