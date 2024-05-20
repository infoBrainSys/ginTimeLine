package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type M map[string]any

const (
	StatusOk                  = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusForbidden           = http.StatusForbidden
	StatusNotFound            = http.StatusNotFound
	StatusInternalServerError = http.StatusInternalServerError
	StatusServiceUnavailable  = http.StatusServiceUnavailable
	StatusGatewayTimeout      = http.StatusGatewayTimeout
)

type response struct {
	Ctx *gin.Context
}

func Response(ctx *gin.Context) *response {
	return &response{
		Ctx: ctx,
	}
}

func (r *response) Ok(m string, d ...*M) {
	r.Ctx.JSON(StatusOk, M{
		"code": StatusOk,
		"msg":  m,
		"data": d,
	})
}

func (r *response) Bad(d *M) {
	r.Ctx.AbortWithStatusJSON(StatusBadRequest, d)
}

func (r *response) ServerError(err error) {
	r.Ctx.AbortWithError(StatusInternalServerError, err)
}

func (r *response) NotFound() {
	r.Ctx.AbortWithStatus(StatusNotFound)
}
