package qa

import (
	provider "github.com/choi006/bbsgo/app/provider/qa"
	"github.com/gohade/hade/framework/gin"
)

type questionDetailParam struct {
	ID int64 `json:"id" binding:"required"`
}

// QuestionDetail 获取问题详情
// @Summary 获取问题详情
// @Description 获取问题详情，包括问题的所有回答
// @Accept  json
// @Produce  json
// @Tags qa
// @Param param body questionDetailParam true "问题id"
// @Success 200 QuestionDTO question "问题详情，带回答和作者"
// @Security ApiKeyAuth
// @Router /question/detail [post]
func (api *QaApi) QuestionDetail(c *gin.Context) {
	// 参数校验
	param := &questionDetailParam{}
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
	if err := qaService.QuestionLoadAuthor(c, question); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}
	if err := qaService.QuestionLoadAnswers(c, question); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}
	if err := qaService.AnswersLoadAuthor(c, &(question.Answers)); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	questionDTO := ConvertQuestionToDTO(question, nil)

	c.ISetOkStatus().IJson(questionDTO)
}
