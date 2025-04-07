package models

import (
	"time"
)

// Post 文章模型
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostStore 文章存储接口
type PostStore interface {
	// FindAll 查找所有文章
	FindAll() ([]*Post, error)

	// FindByID 根据ID查找文章
	FindByID(id int) (*Post, error)

	// Create 创建文章
	Create(post *Post) error

	// Update 更新文章
	Update(post *Post) error

	// Delete 删除文章
	Delete(id int) error
}
