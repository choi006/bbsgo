package qa

import (
	"context"
	"fmt"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type QaService struct {
	container framework.Container //容器
	ormDB     *gorm.DB            //db
	logger    contract.Log        // log
}

func (q *QaService) UpdateQuestion(ctx context.Context, question *Question) error {
	questionDB := &Question{ID: question.ID}
	if err := q.ormDB.WithContext(ctx).First(questionDB).Error; err != nil {
		return errors.WithStack(err)
	}
	
	if question.Title != "" {
		questionDB.Title = question.Title
	}
	if question.Context != "" {
		questionDB.Context = question.Context
	}
	if err := q.ormDB.WithContext(ctx).Save(questionDB).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (q *QaService) DeleteQuestion(ctx context.Context, questionID int64) error {
	questionDB := &Question{ID: questionID}
	if err := q.ormDB.Delete(questionDB).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (q *QaService) GetQuestion(ctx context.Context, questionID int64) (*Question, error) {
	question := &Question{}
	if err := q.ormDB.WithContext(ctx).First(question, questionID).Error; err != nil {
		return nil, err
	}
	return question, nil
}

func (q *QaService) PostQuestion(ctx context.Context, question *Question) error {
	if err := q.ormDB.WithContext(ctx).Create(question).Error; err != nil {
		return err
	}
	return nil
}

func NewQaService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	ormService := container.MustMake(contract.ORMKey).(contract.ORMService)
	logger := container.MustMake(contract.LogKey).(contract.Log)

	db, err := ormService.GetDB()
	if err != nil {
		logger.Error(context.Background(), "获取gromDB错误", map[string]interface{}{
			"err": fmt.Sprintf("%+v", err),
		})
		return nil, err
	}
	return &QaService{container: container, ormDB: db, logger: logger}, nil
}
