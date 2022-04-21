package user

import (
	"github.com/choi006/bbsgo/app/provider/user"
	"github.com/gohade/hade/framework/gin"
)

// Verify godoc
// @Summary 验证注册信息
// @Description 使用token验证注册信息
// @Accept  json
// @Produce json
// @Tags user
// @Param token query string true "注册token"
// @Success 200 {string} Message "注册成功,请进入登录页面"
// @Router /user/register/verify [get]
func (api *UserApi) Verify(c *gin.Context) {
	//验证参数
	token := c.Query("token")
	if token == "" {
		c.ISetStatus(400).IText("参数错误")
		return
	}

	userService := c.MustMake(user.UserKey).(user.Service)
	verified, err := userService.VerifyRegister(c, token)
	if err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	if !verified {
		c.ISetStatus(500).IText("验证错误")
		return
	}

	// 输出
	c.IRedirect("/").IText("注册成功，请进入登录页面")
}
