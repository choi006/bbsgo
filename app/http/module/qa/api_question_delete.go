package qa

import (
	"github.com/choi006/bbsgo/app/http/middleware/auth"
	provider "github.com/choi006/bbsgo/app/provider/qa"
	"github.com/gohade/hade/framework/gin"
)

type questionDeleteParam struct {
	ID int64 `json:"id" binding:"required"`
}

// QuestionDelete 删除问题
// @Summary 删除问题
// @Description 删除问题，同时删除问题中的所有答案
// @Accept  json
// @Produce  json
// @Tags qa
// @Param id query int true "删除id"
// @Success 200 string Msg "操作成功"
// @Security ApiKeyAuth
// @Router /question/delete [get]
func (api *QaApi) QuestionDelete(c *gin.Context) {
	// 参数校验
	param := &questionDeleteParam{}
	if err := c.ShouldBind(param); err != nil {
		c.ISetStatus(400).IText("参数错误")
		return
	}

	qaService := c.MustMake(provider.QaKey).(provider.Service)
	question, err := qaService.GetQuestion(c, param.ID)
	if err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}
	if question == nil {
		c.ISetStatus(500).IText("问题不存在")
		return
	}

	user := auth.GetAuthUser(c)
	if user.ID != question.AuthorID {
		c.ISetStatus(500).IText("无权限操作")
		return
	}

	if err := qaService.DeleteQuestion(c, question.ID); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	c.ISetOkStatus().IText("操作成功")
}
