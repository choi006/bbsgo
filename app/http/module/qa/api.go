package qa

import (
	"github.com/choi006/bbsgo/app/http/middleware/auth"
	"github.com/choi006/bbsgo/app/provider/qa"
	"github.com/gohade/hade/framework/gin"
)

type QaApi struct {
}

func RegisterRoutes(r *gin.Engine) error {
	api := &QaApi{}
	if !r.IsBind(qa.QaKey) {
		r.Bind(&qa.QaProvider{})
	}

	questionApi := r.Group("/question", auth.AuthMiddleware())
	{
		// 创建问题
		questionApi.POST("/create", api.QuestionCreate)
		// 删除问题
		questionApi.POST("/delete", api.QuestionDelete)
		// 更新问题
		questionApi.POST("/edit", api.QuestionEdit)
		// 问题详情
		questionApi.POST("/detail", api.QuestionDetail)
	}
	answerApi := r.Group("/answer", auth.AuthMiddleware())
	{
		// 创建回答
		answerApi.POST("/create", api.AnswerCreate)
		// 删除回答
		answerApi.POST("/delete", api.AnswerDelete)
	}

	return nil
}
