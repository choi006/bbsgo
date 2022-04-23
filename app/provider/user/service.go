package user

import (
	"context"
	"fmt"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

type UserService struct {
	container framework.Container
	logger    contract.Log
	configer  contract.Config
}

func (u *UserService) IsEmailRegister(ctx context.Context, email string) (bool, error) {
	// 判断邮箱是否已经注册了
	ormService := u.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		return true, err
	}
	userDB := &User{}
	if db.Where(&User{Email: email}).First(userDB).Error != gorm.ErrRecordNotFound {
		return true, nil
	}
	return false, nil
}

func (u *UserService) Login(ctx context.Context, user *User) (*User, error) {
	ormService := u.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		return nil, err
	}

	userDB := &User{}
	if err := db.Where("username=?", user.UserName).First(userDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); err != nil {
		return nil, err
	}

	userDB.Password = ""
	return userDB, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genToken(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenNumberValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
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
	// 验证成功将密码存储数据库之前需要加密，不能原文存储进入数据库
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash)

	// 具体在数据库创建用户
	if err := db.Create(user).Error; err != nil {
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

	// 实例化gomail
	d := gomail.NewDialer(host, port, username, password)

	// 组装message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetAddressHeader("To", user.Email, user.UserName)
	m.SetHeader("Subject", "感谢您注册我们的CurryCloud")
	code := GenNumberValidateCode(6)
	cacheService := u.container.MustMake(contract.CacheKey).(contract.CacheService)
	key := fmt.Sprintf("user:register:%s", user.Email)
	if err := cacheService.Set(ctx, key, code, 300*time.Second); err != nil {
		return err
	}
	m.SetBody("text/html", fmt.Sprintf("您的邮件验证码为：%s", code))

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
