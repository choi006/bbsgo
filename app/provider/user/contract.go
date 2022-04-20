package user

import (
	"context"
	"encoding/json"
	"time"
)

const UserKey = "user"

type Service interface {
	// Register 注册用户,注意这里只是将用户注册, 并没有激活, 需要调用
	// 参数：user必填，username，password, email
	// 返回值： user 带上token
	Register(ctx context.Context, user *User) (*User, error)
}

type User struct {
	ID        int64     `gorm:"column:id;primaryKey"` // 代表用户id, 只有注册成功之后才有这个id，唯一表示一个用户
	UserName  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	Email     string    `gorm:"column:email"`
	CreatedAt time.Time `gorm:"column:created_at"`

	Token string `gorm:"-"` // token 可以用作注册token或者登录token
}

// MarshalBinary 实现BinaryMarshaler 接口
func (b *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

// UnmarshalBinary 实现 BinaryUnMarshaler 接口
func (b *User) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, b)
}
