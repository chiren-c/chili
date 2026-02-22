package user

import (
	"github.com/gin-gonic/gin"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 分组注册
	ug := server.Group("/user")
	ug.POST("/signup", c.SignUp)
	ug.POST("/login", c.LoginJWT)
}

func (c *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
}
func (c *UserHandler) LoginJWT(ctx *gin.Context) {}
