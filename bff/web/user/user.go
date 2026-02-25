package user

import (
	jwt3 "github.com/chiren-c/chili/bff/web/jwt"
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/chiren-c/chili/user/domain"
	"github.com/chiren-c/chili/user/errs"
	"github.com/chiren-c/chili/user/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	jwt3.Handler
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 分组注册
	ug := server.Group("/user")
	ug.POST("/signup", c.SignUp)
	ug.POST("/login", c.LoginJWT)
	ug.POST("/refresh_token", c.RefreshToken)
}

// SignUp 用户注册
func (c *UserHandler) SignUp(ctx *gin.Context) {
	type signUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req signUpReq
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

// LoginJWT 用户登录
func (c *UserHandler) LoginJWT(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		c.log.Error("解析请求失败：", loggerx.Error(err))
		return
	}
	u, err := c.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		c.log.Error("执行业务逻辑失败：", loggerx.Error(err))
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	err = c.SetLoginToken(ctx, u.Id)
	if err != nil {
		c.log.Error("执行业务逻辑失败：", loggerx.Error(err))
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "OK"})
}

func (c *UserHandler) RefreshToken(ctx *gin.Context) {
	tokenStr := c.ExtractTokenString(ctx)
	var rc jwt3.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return jwt3.RefreshTokenKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, ginx.Result{Code: errs.UserUnauthorized, Msg: "请登录"})
		return
	}
	// 校验 ssid
	err = c.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// 系统错误或者用户已经主动退出登录了
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = c.SetJWTToken(ctx, rc.Ssid, rc.Id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ginx.Result{Code: errs.UserUnauthorized, Msg: "请登录"})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "刷新成功"})
}

func NewUserHandler(log loggerx.Logger, svc service.UserService,
	jwt jwt3.Handler) *UserHandler {
	return &UserHandler{
		svc:              svc,
		log:              log,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		Handler:          jwt,
	}
}
