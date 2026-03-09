package user

import (
	jwt3 "github.com/chiren-c/chili/bff/web/jwt"
	codeService "github.com/chiren-c/chili/code/service"
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/chiren-c/chili/user/domain"
	"github.com/chiren-c/chili/user/errs"
	"github.com/chiren-c/chili/user/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	phoneRegexPattern    = `/^1[3-9]\d{9}$/`
	bizLogin             = "login"
)

type UserHandler struct {
	svc              service.UserService
	codeSvc          codeService.CodeService
	log              loggerx.Logger
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	phoneRegexExp    *regexp.Regexp
	jwt3.Handler
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 分组注册
	ug := server.Group("/user")
	ug.POST("/signup", c.SignUp)
	ug.POST("/login", c.LoginJWT)
	ug.POST("/login_sms/code/send", c.SendSMSLoginCode)
	ug.POST("/login_sms", c.LoginSMS)
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
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
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
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "注册成功"})
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
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "登录成功"})
}

// RefreshToken 刷新 token
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

// SendSMSLoginCode 发送 sms 登陆 code
func (c *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInvalidInput, Msg: "请输入手机号码"})
	}
	ok, err := c.phoneRegexExp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInvalidInput, Msg: "输入的手机号码有误"})
	}
	err = c.codeSvc.Send(ctx, bizLogin, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "发送成功"})
}

// LoginSMS 登陆
func (c *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := c.phoneRegexExp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInvalidInput, Msg: "输入的手机号码有误"})
	}
	ok, err = c.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		c.log.Error("用户手机号码登录失败", loggerx.Error(err))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInvalidInput, Msg: "验证码错误"})
		return
	}
	u, err := c.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ssid := uuid.New().String()
	err = c.SetJWTToken(ctx, ssid, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.UserInternalServerError, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "登录成功"})
}

func NewUserHandler(log loggerx.Logger, svc service.UserService,
	jwt jwt3.Handler, codeSvc codeService.CodeService) *UserHandler {
	return &UserHandler{
		svc:              svc,
		codeSvc:          codeSvc,
		log:              log,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		phoneRegexExp:    regexp.MustCompile(phoneRegexPattern, regexp.None),
		Handler:          jwt,
	}
}
