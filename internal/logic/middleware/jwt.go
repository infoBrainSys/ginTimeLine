package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	r "timeLineGin/internal/logic/response"
	"timeLineGin/internal/model"
	"timeLineGin/pkg/config"
	"timeLineGin/pkg/redis"
)

const (
	tokenParseError    = "token parse error"
	tokenExpired       = "token expired"
	tokenInvalid       = "token invalid"
	tokenEmpty         = "token empty"
	tokenSignFail      = "token sign fail"
	TokenClaims        = "token"
	TokenValueNotExist = "token value not exist"
)

type JWT struct {
	Ctx *gin.Context
}

// DefaultJwt 初始化 JWT 实例
func DefaultJwt(ctx *gin.Context) *JWT {
	return &JWT{
		Ctx: ctx,
	}
}

// customClaims 自定义声明结构体并内嵌 jwt.StandardClaims
type customClaims struct {
	*jwt.RegisteredClaims
	RemoteIP string
}

// GenerateToken 生成 token
func (j *JWT) GenerateToken(u *model.UserInput) *jwt.Token {
	claims := customClaims{
		&jwt.RegisteredClaims{
			Issuer:    config.GetInstance().Jwt.Iss,
			Subject:   u.Passport,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(gconv.Duration(config.GetInstance().Jwt.Exp) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		j.Ctx.ClientIP(),
	}
	var method *jwt.SigningMethodHMAC
	switch config.GetInstance().Jwt.Method {
	case "SigningMethodHS256":
		method = jwt.SigningMethodHS256
	case "SigningMethodHS384":
		method = jwt.SigningMethodHS384
	case "SigningMethodHS512":
		method = jwt.SigningMethodHS512
	default:
		panic("jwt method error")
	}
	return jwt.NewWithClaims(method, claims)
}

// ValidateToken 校验 tokenString
func (j *JWT) ValidateToken() error {
	token, err := j.ParseToken()
	if err != nil {
		return errors.New(tokenParseError)
	}

	// 校验 tokenString 是否有效
	if !token.Valid {
		return errors.New(tokenInvalid)
	}

	// 是否过期
	t, err := token.Claims.(*customClaims).GetExpirationTime()
	if err != nil || t.Before(time.Now()) {
		return errors.New(tokenExpired)
	}
	return nil
}

// ParseToken 解析并返回 jwt.Token
func (j *JWT) ParseToken() (*jwt.Token, error) {
	return jwt.ParseWithClaims(j.getTokenString(), &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetInstance().Jwt.Secret), nil
	})
}

// ReNewTokenString 刷新 TokenString
func (j *JWT) ReNewTokenString() string {
	token, err := j.ParseToken()
	if err != nil {
		j.Ctx.AbortWithStatusJSON(401, gin.H{
			"code": 401,
			"msg":  "token error",
		})
		return ""
	}

	var expireTime = gconv.Float64(token.Claims.(*customClaims).ExpiresAt)
	switch {
	// 过期时间小于预刷新时间则重签 tokenString
	case expireTime < config.GetInstance().Jwt.Rsn:
		tokenString, _ := j.GenerateToken(&model.UserInput{
			Passport: token.Claims.(*customClaims).Subject,
		}).SignedString(gconv.Bytes(config.GetInstance().Jwt.Secret))
		return tokenString
	// 过期时间大于 预刷新时间则不处理 tokenString
	case expireTime > config.GetInstance().Jwt.Rsn:
		fallthrough
	default:
		return j.getTokenString()
	}
}

// getTokenString 从请求头中获取 tokenString
func (j *JWT) getTokenString() string {
	return j.Ctx.GetHeader("Authorization")[7:]
}

// SetTokenString tokenString 写入到 header 中
func (j *JWT) SetTokenString(tokenString string) {
	j.Ctx.Header("Authorization", "Bearer "+tokenString)
}

// Authorize 身份鉴权中间件
func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 是否获取到 tokenString
		j := DefaultJwt(ctx)
		if tokenStr := j.getTokenString(); tokenStr == "" {
			ctx.AbortWithStatusJSON(401, gin.H{
				"code": 401,
				"msg":  tokenEmpty,
			})
			return
		}

		// 判断 tokenString 是否在黑名单中
		if ok := InBlackLists(ctx); ok == true {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": tokenInvalid})
			return
		}

		// 验证 tokenString 是否有效
		err := j.ValidateToken()
		if err != nil {
			switch err.Error() {
			case tokenExpired:
				ctx.AbortWithStatusJSON(401, gin.H{
					"code": 401,
					"msg":  tokenExpired,
				})
				return
			case tokenInvalid:
				ctx.AbortWithStatusJSON(401, gin.H{
					"code": 401,
					"msg":  tokenInvalid,
				})
				return
			case tokenParseError:
				ctx.AbortWithStatusJSON(401, gin.H{
					"code": 401,
					"msg":  tokenParseError,
				})
				return
			default:
				ctx.AbortWithStatusJSON(401, gin.H{"msg": "unknown err"})
			}
		}

		ctx.Next()

		// 执行刷新 tokenString 策略, 不刷新则写入原来的 tokenString
		j.SetTokenString(j.ReNewTokenString())
	}
}

// SignTokenString 登录签发 token 并写入 header
func SignTokenString() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var u model.UserInput
		if err := ctx.ShouldBind(&u); err != nil {
			r.Response(ctx).Bad(&r.M{"error": err.Error()})
			return
		}

		ctx.Set(TokenClaims, u)
		ctx.Next()

		tokenStr, err := DefaultJwt(ctx).
			GenerateToken(&u).
			SignedString(gconv.Bytes(config.GetInstance().Jwt.Secret))
		if err != nil {
			r.Response(ctx).Bad(&r.M{"error": tokenSignFail})
			return
		}
		ctx.Header("Authorization", "Bearer "+tokenStr)
		ctx.JSON(http.StatusOK, r.M{"code": 200, "msg": "sign in success"})
	}
}

// AddInBlackLists 添加 tokenString 到黑名单
func AddInBlackLists() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		j := DefaultJwt(ctx)
		t, err := j.ParseToken()
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": tokenParseError})
			return
		}

		err = j.ValidateToken()
		if err != nil {
			switch err.Error() {
			case tokenExpired:
				ctx.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": tokenExpired})
				return
			case tokenInvalid:
				ctx.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": tokenInvalid})
				return
			case tokenParseError:
				ctx.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": tokenParseError})
				return
			}
		}

		sub, err := t.Claims.GetSubject()
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": TokenValueNotExist})
			return
		}

		expTime, _ := t.Claims.GetExpirationTime()
		if !expTime.After(time.Now()) {
			ctx.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": tokenExpired})
			return
		}
		redis.GetInstance().Set(
			ctx,
			sub+"-ExpiredToken",        // 构造格式 [sub]-ExpiredToken
			sub+"-"+j.getTokenString(), // 构造格式 [sub]-[token]
			expTime.Sub(time.Now()),
		)
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok"})
	}
}

func blackListKey(sub string) string {
	return sub + "-ExpiredToken"
}

// InBlackLists 判断 tokenString 是否在黑名单中
func InBlackLists(ctx *gin.Context) bool {
	j := DefaultJwt(ctx)

	err := j.ValidateToken()
	if err != nil {
		return false
	}

	t, err := j.ParseToken()
	if err != nil {
		return false
	}

	sub, _ := t.Claims.GetSubject()
	if err := redis.GetInstance().Get(ctx, blackListKey(sub)).Err(); err != nil {
		return false
	}
	return true
}
