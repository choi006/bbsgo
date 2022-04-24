package auth

import (
	"github.com/choi006/bbsgo/app/provider/user"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/gin"
	"github.com/golang-jwt/jwt"
	"time"
)

type MyCustomClaims struct {
	ID       int64
	Username string
	jwt.StandardClaims
}

// AuthMiddleware 登录中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.ISetStatus(401).IText("请登录后操作")
			c.Abort()
			return
		}
		claim, err := ParseToken(c, token)
		if err != nil || claim == nil {
			c.ISetStatus(401).IText("请登录后操作")
			c.Abort()
			return
		}

		authUser := &user.User{
			ID:       claim.ID,
			UserName: claim.Username,
		}

		c.Set("auth_user", authUser)

		c.Next()
	}
}

// GetAuthUser 获取已经验证的用户
func GetAuthUser(c *gin.Context) *user.User {
	t, exist := c.Get("auth_user")
	if !exist {
		return nil
	}
	return t.(*user.User)
}

// GenerateToken 生成jwt
func GenerateToken(c *gin.Context, user *user.User) (string, error) {
	configer := c.MustMake(contract.ConfigKey).(contract.Config)
	// 配置服务中生成jwt需要的参数
	secretKey := configer.GetString("app.jwt.secret_key")
	issuer := configer.GetString("app.jwt.issuer")
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := MyCustomClaims{
		ID:       user.ID,
		Username: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	return token, err
}

// ParseToken 解析jwt
func ParseToken(c *gin.Context, token string) (*MyCustomClaims, error) {
	configer := c.MustMake(contract.ConfigKey).(contract.Config)
	// 配置服务中生成jwt需要的参数
	secretKey := configer.GetString("app.jwt.secret_key")
	tokenClaims, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*MyCustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
