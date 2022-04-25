package qa

import (
	"context"
	"github.com/choi006/bbsgo/app/provider/user"
	"gorm.io/gorm"
	"time"
)

const QaKey = "qa"

type Service interface {
	// GetQuestion 获取某个问题详情，question简化结构
	GetQuestion(ctx context.Context, questionID int64) (*Question, error)
	// PostQuestion 上传某个问题
	// ctx必须带上操作人id
	PostQuestion(ctx context.Context, question *Question) error
	// DeleteQuestion 删除问题，同时删除对应的回答
	// ctx必须带操作人信息
	DeleteQuestion(ctx context.Context, questionID int64) error
	// UpdateQuestion 代表更新问题, 只会对比其中的context，title两个字段，其他字段不会对比
	// ctx必须带操作人
	UpdateQuestion(ctx context.Context, question *Question) error
}

// Question 代表问题
type Question struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Title     string         `gorm:"column:title;comment:标题"`
	Context   string         `gorm:"column:context;comment:内容"`
	AuthorID  int64          `gorm:"column:author_id;comment:作者id;not null;default:0"`
	AnswerNum int            `gorm:"column:answer_num;comment:回答数;not null;default:0"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime;autoCreateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Author    *user.User     `gorm:"foreignKey:AuthorID"`
	Answers   []*Answer      `gorm:"foreignKey:QuestionID"`
}

// Answer 代表一个回答
type Answer struct {
	ID         int64          `gorm:"column:id;primaryKey"`
	QuestionID int64          `gorm:"column:question_id;index;comment:问题id;not null;default 0"`
	Context    string         `gorm:"column:context;comment:内容"`
	AuthorID   int64          `gorm:"column:author_id;comment:作者id;not null;default:0"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime;autoCreateTime;<-:false;comment:更新时间"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Author     *user.User     `gorm:"foreignKey:AuthorID"`
	Question   *Question      `gorm:"foreignKey:QuestionID"`
}
