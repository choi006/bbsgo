package user

import (
	"context"
	"encoding/json"
	"time"
)

const UserKey = "user"

// Service 用户相关的服务
type Service interface {
	// Register 注册用户,注意这里只是将用户注册, 并没有激活, 需要调用
	// 参数：user必填，username，password, email
	// 返回值： user 带上token
	Register(ctx context.Context, user *User) (*User, error)
	// IsEmailRegister 用户邮箱是否注册
	// 参数：user必填： email
	IsEmailRegister(ctx context.Context, email string) (bool, error)
	// SendRegisterMail 发送注册的邮件
	// 参数：user必填： username, password, email, token
	SendRegisterMail(ctx context.Context, user *User) error
	// Login 登录相关，使用用户名密码登录，获取完成User信息
	Login(ctx context.Context, user *User) (*User, error)
}

// User 代表一个用户，注意这里的用户信息字段在不同接口和参数可能为空
type User struct {
	ID        int64     `gorm:"column:id;primaryKey"` // 代表用户id, 只有注册成功之后才有这个id，唯一表示一个用户
	UserName  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	Email     string    `gorm:"column:email"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

// MarshalBinary 实现BinaryMarshaler 接口
func (b *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

// UnmarshalBinary 实现 BinaryUnMarshaler 接口
func (b *User) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, b)
}
