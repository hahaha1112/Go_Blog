package controllers

import (
	"goblog/db"
	"goblog/models"
	"goblog/utils"
	"html/template"
	"net/http"
	"time"
)

// LoginFormHandler 处理登录表单请求
func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	// 检查用户是否已登录
	user := utils.GetUserFromSession(r)
	if user != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 渲染模板
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/users/login.html",
	)
	if err != nil {
		http.Error(w, "模板解析错误", http.StatusInternalServerError)
		return
	}

	// 获取当前年份
	currentYear := time.Now().Year()

	data := map[string]interface{}{
		"Title":       "用户登录",
		"CurrentYear": currentYear,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "模板渲染错误", http.StatusInternalServerError)
	}
}

// LoginProcessHandler 处理登录请求
func LoginProcessHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "表单解析错误", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	username := r.FormValue("username")
	password := r.FormValue("password")

	// 简单验证
	if username == "" || password == "" {
		http.Error(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}

	// 获取存储实例
	store, err := db.NewSQLiteStore("./goblog.db")
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	defer store.Close()

	// 认证用户
	user, err := store.Authenticate(username, password)
	if err != nil {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	// 设置会话
	if err := utils.SetUserSession(w, r, user); err != nil {
		http.Error(w, "无法创建会话", http.StatusInternalServerError)
		return
	}

	// 重定向到首页
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutHandler 处理登出请求
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// 清除会话
	utils.ClearUserSession(w, r)

	// 重定向到首页
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// RegisterFormHandler 处理注册表单请求
func RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	// 检查用户是否已登录
	user := utils.GetUserFromSession(r)
	if user != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 渲染模板
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/users/register.html",
	)
	if err != nil {
		http.Error(w, "模板解析错误", http.StatusInternalServerError)
		return
	}

	// 获取当前年份
	currentYear := time.Now().Year()

	data := map[string]interface{}{
		"Title":       "用户注册",
		"CurrentYear": currentYear,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "模板渲染错误", http.StatusInternalServerError)
	}
}

// RegisterProcessHandler 处理注册请求
func RegisterProcessHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "表单解析错误", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// 简单验证
	if username == "" || email == "" || password == "" {
		http.Error(w, "所有字段都必须填写", http.StatusBadRequest)
		return
	}

	if password != confirmPassword {
		http.Error(w, "两次密码输入不一致", http.StatusBadRequest)
		return
	}

	// 创建用户
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	// 获取存储实例
	store, err := db.NewSQLiteStore("./goblog.db")
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	defer store.Close()

	// 检查用户名是否已存在
	_, err = store.FindUserByUsername(username)
	if err == nil {
		http.Error(w, "用户名已存在", http.StatusConflict)
		return
	}

	// 检查邮箱是否已存在
	_, err = store.FindUserByEmail(email)
	if err == nil {
		http.Error(w, "邮箱已存在", http.StatusConflict)
		return
	}

	// 保存用户
	if err := store.CreateUser(user); err != nil {
		http.Error(w, "无法创建用户", http.StatusInternalServerError)
		return
	}

	// 设置会话
	if err := utils.SetUserSession(w, r, user); err != nil {
		http.Error(w, "无法创建会话", http.StatusInternalServerError)
		return
	}

	// 重定向到首页
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
