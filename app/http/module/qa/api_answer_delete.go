package qa

import (
	"github.com/choi006/bbsgo/app/http/middleware/auth"
	provider "github.com/choi006/bbsgo/app/provider/qa"
	"github.com/gohade/hade/framework/gin"
)

type answerDeleteParam struct {
	ID int64 `json:"id" binding:"required"`
}

// AnswerDelete 代表删除回答
// @Summary 删除回答
// @Description 删除回答
// @Accept  json
// @Produce  json
// @Tags qa
// @Param answerDeleteParam body answerDeleteParam true "删除id"
// @Success 200 string Msg "操作成功"
// @Security ApiKeyAuth
// @Router /answer/delete [post]
func (api *QaApi) AnswerDelete(c *gin.Context) {
	// 参数校验
	param := &answerDeleteParam{}
	if err := c.ShouldBind(param); err != nil {
		c.ISetStatus(400).IText("参数错误")
		return
	}

	user := auth.GetAuthUser(c)
	if user == nil {
		c.ISetStatus(500).IText("请登录后再操作")
		return
	}

	qaService := c.MustMake(provider.QaKey).(provider.Service)
	answer, err := qaService.GetAnswer(c, param.ID)
	if err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}
	if answer.AuthorID != user.ID {
		c.ISetStatus(500).IText("没有权限作此操作")
		return
	}

	if err := qaService.DeleteAnswer(c, answer); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	c.ISetOkStatus().IText("操作成功")
}
