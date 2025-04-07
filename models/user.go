package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // 不输出到JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserStore 用户存储接口
type UserStore interface {
	// FindByID 根据ID查找用户
	FindByID(id int) (*User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(email string) (*User, error)

	// Create 创建用户
	Create(user *User) error

	// Update 更新用户
	Update(user *User) error

	// Delete 删除用户
	Delete(id int) error

	// Authenticate 认证用户
	Authenticate(username, password string) (*User, error)
}
