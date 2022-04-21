package user

import (
	"context"
	"fmt"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
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
	// 判断邮箱是否已经注册了
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

func (u *UserService) SendRegisterMail(ctx context.Context, user *User) error {
	logger := u.container.MustMake(contract.LogKey).(contract.Log)
	configer := u.container.MustMake(contract.ConfigKey).(contract.Config)

	// 配置服务中获取发送邮件需要的参数
	host := configer.GetString("app.smtp.host")
	port := configer.GetInt("app.smtp.port")
	username := configer.GetString("app.smtp.username")
	password := configer.GetString("app.smtp.password")
	from := configer.GetString("app.smtp.from")
	domain := configer.GetString("app.domain")

	// 实例化gomail
	d := gomail.NewDialer(host, port, username, password)

	// 组装message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetAddressHeader("To", user.Email, user.UserName)
	m.SetHeader("Subject", "感谢您注册我们的CurryCloud")
	link := fmt.Sprintf("%v/user/register/verify?token=%v", domain, user.Token)
	m.SetBody("text/html", fmt.Sprintf("请点击下面的链接完成注册：%s", link))

	// 发送电子邮件
	if err := d.DialAndSend(m); err != nil {
		logger.Error(ctx, "send email error", map[string]interface{}{
			"err":     err,
			"message": m,
		})
		return err
	}
	return nil
}

func NewUserService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	logger := container.MustMake(contract.LogKey).(contract.Log)
	configer := container.MustMake(contract.ConfigKey).(contract.Config)
	return &UserService{container: container, logger: logger, configer: configer}, nil
}
