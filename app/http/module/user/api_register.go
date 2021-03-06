package user

import (
	"fmt"
	"github.com/choi006/bbsgo/app/http/middleware/auth"
	provider "github.com/choi006/bbsgo/app/provider/user"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/gin"
	"time"
)

type registerParam struct {
	UserName     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required,gte=6"`
	Email        string `json:"email" binding:"required,gte=6"`
	ValidateCode string `json:"validate_code" binding:"required,gte=6"`
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口
// @Accept  json
// @Produce json
// @Tags user
// @Param registerParam body registerParam true "注册参数"
// @Success 200 string Message "token"
// @Router /user/register [post]
func (api *UserApi) Register(c *gin.Context) {
	// 验证参数
	userService := c.MustMake(provider.UserKey).(provider.Service)
	logger := c.MustMake(contract.LogKey).(contract.Log)

	param := &registerParam{}
	if err := c.ShouldBind(param); err != nil {
		c.ISetStatus(400).IText("参数错误")
		return
	}

	// 验证邮箱验证码
	cacheService := c.MustMake(contract.CacheKey).(contract.CacheService)
	key := fmt.Sprintf("user:register:%s", param.Email)
	code, err := cacheService.Get(c, key)
	if err != nil || param.ValidateCode != code {
		c.ISetStatus(500).IText("验证码错误")
		return
	}

	// 登录
	model := &provider.User{
		UserName:  param.UserName,
		Password:  param.Password,
		Email:     param.Email,
		CreatedAt: time.Now(),
	}
	// 注册
	user, err := userService.Register(c, model)
	if err != nil {
		logger.Error(c, err.Error(), map[string]interface{}{
			"stack": fmt.Sprintf("%+v", err),
		})
		c.ISetStatus(500).IText(err.Error())
		return
	}
	if user == nil {
		c.ISetStatus(500).IText("注册失败")
		return
	}

	token, err := auth.GenerateToken(c, user)
	if err != nil {
		c.ISetStatus(500).IText("生成token失败")
		return
	}

	// 输出
	c.ISetOkStatus().IText(token)
	return
}
