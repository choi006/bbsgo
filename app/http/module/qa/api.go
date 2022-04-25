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
	}

	return nil
}
