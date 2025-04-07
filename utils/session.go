package utils

import (
	"encoding/json"
	"goblog/models"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

var (
	// store 是会话存储
	store = sessions.NewCookieStore([]byte("goblog-secret-key"))

	// sessionName 是会话的名称
	sessionName = "goblog-session"
)

// SetUserSession 设置用户会话
func SetUserSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	// 获取会话
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}

	// 创建安全的用户数据（不包含密码）
	safeUser := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}

	// 序列化用户数据
	userData, err := json.Marshal(safeUser)
	if err != nil {
		return err
	}

	// 设置会话数据
	session.Values["user"] = userData
	session.Options.MaxAge = 86400 * 7 // 7天

	// 保存会话
	return session.Save(r, w)
}

// GetUserFromSession 从会话中获取用户
func GetUserFromSession(r *http.Request) *models.User {
	// 获取会话
	session, err := store.Get(r, sessionName)
	if err != nil {
		return nil
	}

	// 检查用户数据是否存在
	userData, ok := session.Values["user"]
	if !ok {
		return nil
	}

	// 反序列化用户数据
	var userMap map[string]interface{}
	if err := json.Unmarshal(userData.([]byte), &userMap); err != nil {
		return nil
	}

	// 创建用户对象
	user := &models.User{
		ID:        int(userMap["id"].(float64)),
		Username:  userMap["username"].(string),
		Email:     userMap["email"].(string),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	return user
}

// ClearUserSession 清除用户会话
func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	// 获取会话
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}

	// 删除会话数据
	delete(session.Values, "user")
	session.Options.MaxAge = -1 // 立即过期

	// 保存会话
	return session.Save(r, w)
}
