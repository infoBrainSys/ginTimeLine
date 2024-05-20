package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/util/gconv"
	b "timeLineGin/internal/logic/business"
	"timeLineGin/internal/logic/middleware"
	r "timeLineGin/internal/logic/response"
	"timeLineGin/internal/model"
	"timeLineGin/internal/service"
	"timeLineGin/pkg/config"
)

const (
	SignUpSuccess = "注册成功"
	userExist     = "用户已存在"
)

func SignUp(ctx *gin.Context) {
	bis := b.Ctx(ctx)
	var user model.UserInput
	if err := bis.Ctx.ShouldBind(&user); err != nil {
		r.Response(bis.Ctx).Bad(&r.M{"msg": err.Error()})
		return
	}

	if err := service.Sign().SignUp(&user); err != nil {
		var e *mysql.MySQLError
		errors.As(err, &e)
		switch e.Number {
		case 1062:
			r.Response(bis.Ctx).Bad(&r.M{"msg": userExist})
			return
		}
	}

	tokenStr, _ := middleware.DefaultJwt(bis.Ctx).
		GenerateToken(&user).
		SignedString(gconv.Bytes(config.GetInstance().Jwt.Secret))
	bis.Ctx.Header("Authorization", "Bearer "+tokenStr)
	r.Response(bis.Ctx).Ok(SignUpSuccess, nil)
}
