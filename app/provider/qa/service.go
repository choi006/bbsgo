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

func (q *QaService) QuestionsLoadAuthor(ctx context.Context, questions *[]*Question) error {
	if err := q.ormDB.WithContext(ctx).Preload("Author").Find(questions).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (q *QaService) GetQuestions(ctx context.Context, paper *Pager) ([]*Question, error) {
	questions := make([]*Question, 0, paper.Size)
	total := int64(0)
	if err := q.ormDB.WithContext(ctx).Model(&Question{}).Count(&total).Error; err != nil {
		paper.Total = total
	}
	if err := q.ormDB.WithContext(ctx).Order("created_at desc").Offset(paper.Start).Limit(paper.Size).Find(&questions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*Question{}, nil
		}
		return nil, err
	}
	return questions, nil
}

func (q *QaService) AnswersLoadAuthor(ctx context.Context, answers *[]*Answer) error {
	if answers == nil {
		return nil
	}
	ids := make([]int64, 0)
	for _, answer := range *answers {
		ids = append(ids, answer.ID)
	}
	if len(ids) == 0 {
		return nil
	}
	if err := q.ormDB.WithContext(ctx).Preload("Author").Order("created_at desc").Find(answers, ids).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (q *QaService) QuestionLoadAnswers(ctx context.Context, question *Question) error {
	if err := q.ormDB.WithContext(ctx).Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("answers.created_at desc")
	}).First(question).Error; err != nil {
		return err
	}
	return nil
}

func (q *QaService) QuestionLoadAuthor(ctx context.Context, question *Question) error {
	if err := q.ormDB.WithContext(ctx).Preload("Author").First(question).Error; err != nil {
		return err
	}
	return nil
}

func (q *QaService) DeleteAnswer(ctx context.Context, answerID int64) error {
	answerDB := &Answer{ID: answerID}
	if err := q.ormDB.WithContext(ctx).Delete(answerDB).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (q *QaService) GetAnswer(ctx context.Context, answerId int64) (*Answer, error) {
	answer := &Answer{}
	if err := q.ormDB.WithContext(ctx).First(answer, answerId).Error; err != nil {
		return nil, err
	}
	return answer, nil
}

func (q *QaService) PostAnswer(ctx context.Context, answer *Answer) error {
	if answer.QuestionID == 0 {
		return errors.New("问题不存在")
	}
	// 必须使用事务
	err := q.ormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		question := &Question{ID: answer.QuestionID}
		// 获取问题
		if err := tx.First(question).Error; err != nil {
			return err
		}
		// 增加回答
		if err := tx.Create(answer).Error; err != nil {
			return err
		}
		// 问题回答数量+1
		question.AnswerNum = question.AnswerNum + 1
		if err := tx.Save(question).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
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
