package business

import (
	"github.com/gin-gonic/gin"
	"timeLineGin/internal/logic/middleware"
)

type M map[string]any

type BaseContext struct {
	Ctx   *gin.Context
	Jwt   *middleware.JWT
	Token string
}

func Ctx(ctx *gin.Context) *BaseContext {
	return &BaseContext{
		Ctx: ctx,
		Jwt: middleware.DefaultJwt(ctx),
	}
}
