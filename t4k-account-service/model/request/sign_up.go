package request

type SignUp struct {
	Name     string `form:"username" binding:"required,max=10,min=3"`
	Password string `form:"password" binding:"required,max=30,min=5"`
}
