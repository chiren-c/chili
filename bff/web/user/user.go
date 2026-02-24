package user

import (
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/chiren-c/chili/user/domain"
	"github.com/chiren-c/chili/user/errs"
	"github.com/chiren-c/chili/user/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	svc              service.UserService
	log              loggerx.Logger
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 分组注册
	ug := server.Group("/user")
	ug.POST("/signup", c.SignUp)
	ug.POST("/login", c.LoginJWT)
}

// SignUp 用户注册接口
func (c *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		c.log.Error("解析请求失败：", loggerx.Error(err))
		return
	}
	ok, err := c.emailRegexExp.MatchString(req.Email)
	if err != nil {
		c.log.Error("执行业务逻辑失败：", loggerx.Error(err))
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInvalidInput, Msg: "邮箱错误"})
		return
	}
	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInvalidInput, Msg: "邮箱密码不能为空"})
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInvalidInput, Msg: "输入两次密码不相等"})
		return
	}

	ok, err = c.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		c.log.Error("执行业务逻辑失败：", loggerx.Error(err))
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInvalidInput, Msg: "密码必须包含数字、特殊字符，并且长度不能小于 8 位"})
		return
	}
	err = c.svc.Signup(ctx, domain.User{Email: req.Email, Password: req.Password})
	if err != nil {
		c.log.Error("执行业务逻辑失败：", loggerx.Error(err))
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "OK"})
}
func (c *UserHandler) LoginJWT(ctx *gin.Context) {}

func NewUserHandler(log loggerx.Logger, svc service.UserService) *UserHandler {
	return &UserHandler{
		svc:              svc,
		log:              log,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}
