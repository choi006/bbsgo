package qa

import (
    "fmt"
    "github.com/choi006/bbsgo/app/http/middleware/auth"
    provider "github.com/choi006/bbsgo/app/provider/qa"
    "github.com/gohade/hade/framework/gin"
)

type questionCreateParam struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content" binding:"required"`
}

// QuestionCreate 代表创建问题
// @Summary 创建问题
// @Description 创建问题
// @Accept  json
// @Produce  json
// @Tags provider
// @Param questionCreateParam body questionCreateParam true "创建问题参数"
// @Success 200 string Msg "操作成功"
// @Security ApiKeyAuth
// @Router /question/create [post]
func (api *QaApi) QuestionCreate(c *gin.Context) {
    fmt.Println("QuestionCreate")
    // 参数校验
    param := &questionCreateParam{}
    if err := c.ShouldBind(param); err != nil {
        c.ISetStatus(400).IText("参数错误")
        return
    }

    user := auth.GetAuthUser(c)
    if user == nil {
        c.ISetStatus(500).IText("无权限操作")
        return
    }

    question := &provider.Question{
        Title:    param.Title,
        Context:  param.Content,
        AuthorID: user.ID,
    }
    qaService := c.MustMake(provider.QaKey).(provider.Service)
    if err := qaService.PostQuestion(c, question); err != nil {
        c.ISetStatus(500).IText(err.Error())
        return
    }

    c.ISetOkStatus().IText("操作成功")
}
