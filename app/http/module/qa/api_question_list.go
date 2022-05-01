package qa

import (
	"fmt"
	provider "github.com/choi006/bbsgo/app/provider/qa"
	"github.com/gohade/hade/framework/gin"
)

type questionListParam struct {
	Start int `json:"start"`
	Size  int `json:"size" binding:"required"`
}

// QuestionList 获取问题列表
// @Summary 获取问题列表
// @Description 获取问题列表，包含作者信息，不包含回答
// @Accept  json
// @Produce  json
// @Tags qa
// @Param questionListParam body questionListParam true "获取问题列表参数"
// @Success 200 {array} QuestionDTO questions "问题列表"
// @Security ApiKeyAuth
// @Router /question/list [post]
func (api *QaApi) QuestionList(c *gin.Context) {
	// 参数校验
	param := &questionListParam{}
	if err := c.ShouldBind(param); err != nil {
		c.ISetStatus(400).IText("参数错误")
		return
	}
	fmt.Printf("%+v", param)
	paper := &provider.Pager{
		Start: param.Start,
		Size:  param.Size,
	}
	logger := c.MustMakeLog()
	logger.Debug(c, "get param", map[string]interface{}{
		"paper": paper,
	})
	qaService := c.MustMake(provider.QaKey).(provider.Service)
	questions, err := qaService.GetQuestions(c, paper)
	if err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}
	if len(questions) == 0 {
		c.ISetOkStatus().IJson([]*QuestionDTO{})
		return
	}
	if err := qaService.QuestionsLoadAuthor(c, &questions); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}
	questionsDTO := ConvertQuestionsToDTO(questions)

	c.ISetOkStatus().IJson(questionsDTO)
}
