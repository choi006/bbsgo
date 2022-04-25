package qa

import (
	"github.com/choi006/bbsgo/app/http/middleware/auth"
	provider "github.com/choi006/bbsgo/app/provider/qa"
	"github.com/gohade/hade/framework/gin"
)

type questionEditParam struct {
	ID      int64  `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// QuestionEdit 编辑问题
// @Summary 编辑问题
// @Description 编辑问题
// @Accept  json
// @Produce  json
// @Tags qa
// @Param questionEditParam body questionEditParam true "编辑问题参数"
// @Success 200 string Msg "操作成功"
// @Router /question/edit [post]
func (api *QaApi) QuestionEdit(c *gin.Context) {
	// 参数校验
	param := &questionEditParam{}
	if err := c.ShouldBind(param); err != nil {
		c.ISetStatus(400).IText("参数错误")
		return
	}

	qaService := c.MustMake(provider.QaKey).(provider.Service)
	questionOld, err := qaService.GetQuestion(c, param.ID)
	if err != nil || questionOld == nil {
		c.ISetStatus(500).IText("操作的问题不存在")
		return
	}

	authUser := auth.GetAuthUser(c)
	if authUser == nil || authUser.ID != questionOld.AuthorID {
		c.ISetStatus(500).IText("无权限操作")
		return
	}

	question := &provider.Question{
		ID:      param.ID,
		Title:   param.Title,
		Context: param.Content,
	}
	if err := qaService.UpdateQuestion(c, question); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	c.ISetOkStatus().IText("操作成功")
}
