package user

import (
	provider "github.com/choi006/bbsgo/app/provider/user"
	"github.com/gohade/hade/framework/gin"
)

type ValidateCodeGetParam struct {
	Email string `json:"email" binding:"required,gte=6"`
}

// PreRegister godoc
// @Summary 用户预注册，获取邮箱验证码
// @Description 用户预注册接口
// @Accept  json
// @Produce json
// @Tags user
// @Param ValidateCodeGetParam body ValidateCodeGetParam true "预注册参数"
// @Success 200 string Message "预注册成功"
// @Router /user/preregister [get]
func (api *UserApi) PreRegister(c *gin.Context) {
	// 验证参数
	param := &ValidateCodeGetParam{}
	if err := c.ShouldBind(param); err != nil {
		c.ISetStatus(400).IText("参数错误")
		return
	}

	model := &provider.User{
		Email: param.Email,
	}

	userService := c.MustMake(provider.UserKey).(provider.Service)
	isRegister, err := userService.IsEmailRegister(c, param.Email)
	if err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}
	if isRegister {
		c.ISetStatus(500).IText("该邮箱已注册")
		return
	}
	// 发送邮件验证码
	if err := userService.SendRegisterMail(c, model); err != nil {
		c.ISetStatus(500).IText("发送邮件验证码失败")
		return
	}

	c.ISetOkStatus().IText("发送邮件验证码成功，请前往邮箱查看邮件")
	return
}
