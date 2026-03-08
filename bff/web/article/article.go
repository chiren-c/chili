package article

import (
	"github.com/chiren-c/chili/article/domain"
	"github.com/chiren-c/chili/article/errs"
	"github.com/chiren-c/chili/article/service"
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ArticleHandler struct {
	log loggerx.Logger
	svc service.ArticleService
}

func (a *ArticleHandler) RegisterRoutes(s *gin.Engine) {
	ag := s.Group("/article")
	ag.POST("/list", a.list)
	ag.POST("/save", a.save)
	ag.POST("/detail", a.detail)
	ag.POST("/publish", a.publish)
}

func (a *ArticleHandler) list(ctx *gin.Context) {
	var usr ginx.UserClaims
	user, ok := ctx.Get("user")
	if !ok {
		a.log.Error("无法获得 claims：", loggerx.String("path", ctx.Request.URL.Path))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	usr, ok = user.(ginx.UserClaims)
	if !ok {
		a.log.Error("无法获得 claims：", loggerx.String("path", ctx.Request.URL.Path))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	type listReq struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}
	var req listReq
	if err := ctx.Bind(&req); err != nil {
		a.log.Error("解析请求失败：", loggerx.Error(err))
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.ArticleInternalServerError, Msg: "系统错误"})
		return
	}
	if req.Limit > 100 {
		a.log.Error("获取列表信息失败，LIMIT过大")
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.ArticleInternalServerError, Msg: "系统错误"})
		return
	}
	arts, err := a.svc.List(ctx, usr.Id, req.Limit, req.Offset)
	if err != nil {
		a.log.Error("获取列表信息失败", loggerx.Error(err))
		ctx.JSON(http.StatusOK, ginx.Result{Code: errs.ArticleInternalServerError, Msg: "系统错误"})
		return
	}
	data := slice.Map[domain.ArticleAuthor, ArticleAuthorVo](arts,
		func(idx int, src domain.ArticleAuthor) ArticleAuthorVo {
			return ArticleAuthorVo{
				Id:         src.Id,
				Title:      src.Title,
				Abstract:   src.Abstract(),
				Status:     src.Status.ToUint8(),
				StatusText: src.Status.ToString(),
				Ctime:      src.Ctime.Format(time.DateTime),
				Utime:      src.Utime.Format(time.DateTime),
			}
		})
	ctx.JSON(http.StatusOK, ginx.Result{Data: data})
}

func (a *ArticleHandler) save(ctx *gin.Context) {

}

func (a *ArticleHandler) detail(ctx *gin.Context) {

}

func (a *ArticleHandler) publish(ctx *gin.Context) {

}

func NewArticleHandler(log loggerx.Logger, svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{log: log, svc: svc}
}
