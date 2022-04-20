package user

import (
	"context"
	"fmt"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type UserService struct {
	container framework.Container
	logger    contract.Log
	configer  contract.Config
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genToken(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (u *UserService) Register(ctx context.Context, user *User) (*User, error) {
	ormService := u.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		return nil, err
	}
	userDB := &User{}
	if db.Where(&User{Email: user.Email}).First(userDB).Error != gorm.ErrRecordNotFound {
		return nil, errors.New("邮箱已注册用户，不能重复注册")
	}
	if db.Where(&User{UserName: user.UserName}).First(userDB).Error != gorm.ErrRecordNotFound {
		return nil, errors.New("用户名已经被注册，请换一个用户名")
	}

	token := genToken(10)
	user.Token = token

	// 将请求注册进入redis，保存一天
	cacheService := u.container.MustMake(contract.CacheKey).(contract.CacheService)

	key := fmt.Sprintf("user:register:%v", user.Token)
	if err := cacheService.SetObj(ctx, key, user, 24*time.Hour); err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	logger := container.MustMake(contract.LogKey).(contract.Log)
	configer := container.MustMake(contract.ConfigKey).(contract.Config)
	return &UserService{container: container, logger: logger, configer: configer}, nil
}
