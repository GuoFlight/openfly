package v1

import (
	"github.com/GuoFlight/gerror"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"openfly/conf"
	"openfly/logger"
	"strings"
	"time"
)

const HeaderNameAuth = "Authorization"

type JwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type ReqLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetTokenFromReq(c)
		if token == "" {
			c.JSON(401, Res{
				Code: 1,
				Msg:  "no token",
			})
			c.Abort()
			return
		}
		_, gerr := ParseToken(token)
		if gerr != nil {
			c.JSON(401, Res{
				Code: 1,
				Msg:  gerr.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetTokenFromReq(c *gin.Context) string {
	for k, v := range c.Request.Header {
		if k == HeaderNameAuth {
			if len(v) == 0 {
				return ""
			}
			return strings.Trim(v[0], `"`)
		}
	}
	return ""
}
func Login(c *gin.Context) {
	var req ReqLogin
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, Res{
			Code: 400,
			Msg:  "invalid request",
		})
		return
	}
	if req.Username == conf.GConf.API.AdminUser && req.Password == conf.GConf.API.AdminPassword {
		token, err := GenToken(conf.GConf.API.AdminUser, conf.GConf.API.AdminPassword)
		if err != nil {
			c.JSON(500, Res{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}
		c.JSON(200, Res{
			Code: 0,
			Data: token,
		})
		return
	} else {
		c.JSON(400, Res{
			Code: 400,
			Msg:  "invalid username or password",
		})
		return
	}
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*JwtClaims, *gerror.Gerr) {
	logger.GLogger.Tracef("即将解析Token:%s", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(conf.GConf.API.AdminPassword), nil
	})
	if err != nil {
		return nil, logger.PrintErr(gerror.NewErr(err.Error()), nil)
	}
	if !token.Valid { // 过期也会返回error
		return nil, gerror.NewErr("invalid token")
	}

	if claims, ok := token.Claims.(*JwtClaims); ok { // 校验token
		return claims, nil
	}
	return nil, gerror.NewErr("invalid token")
}

func GenToken(username, password string) (string, error) {
	c := JwtClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(conf.GConf.API.Expire) * time.Second).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(password))
}
