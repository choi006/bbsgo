package user

import (
	"github.com/choi006/bbsgo/app/provider/user"
	"github.com/gohade/hade/framework/gin"
)

type UserApi struct {
}

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine) error {
	api := &UserApi{}
	if !r.IsBind(user.UserKey) {
		r.Bind(&user.UserProvider{})
	}

	// 预注册
	r.GET("/user/preregister", api.PreRegister)
	// 注册
	r.POST("/user/register", api.Register)
	// 登录
	r.POST("/user/login", api.Login)
	// 登出
	r.GET("/user/logout", api.Logout)

	return nil
}
