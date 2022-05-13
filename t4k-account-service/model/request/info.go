package request

type Info struct {
	UserId uint64 `form:"user_id" binding:"required"`
}
