package controller

import (
	"github.com/gin-gonic/gin"
	b "timeLineGin/internal/logic/business"
	"timeLineGin/internal/logic/middleware"
	r "timeLineGin/internal/logic/response"
	"timeLineGin/internal/model"
	"timeLineGin/internal/service"
)

const (
	errKey              = "error"
	SignInSuccess       = "SignIn success"
	tokenClaimsNotExist = "token claims not exist"
)

func SignIn(ctx *gin.Context) {
	bis := b.Ctx(ctx)
	u, exist := ctx.Get(middleware.TokenClaims)
	if !exist {
		r.Response(bis.Ctx).Bad(&r.M{errKey: tokenClaimsNotExist})
	}
	m := u.(model.UserInput)
	if err := service.Sign().SignIn(&m); err != nil {
		r.Response(bis.Ctx).Bad(&r.M{errKey: err.Error()})
		return
	}
}
