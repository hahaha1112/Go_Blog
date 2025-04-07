package controllers

import (
	"goblog/db"
	"goblog/models"
	"goblog/utils"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ListPostsHandler 处理文章列表请求
func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	// 获取存储实例
	store, err := db.NewSQLiteStore("./goblog.db")
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	defer store.Close()

	// 获取所有文章
	posts, err := store.FindAllPosts()
	if err != nil {
		http.Error(w, "无法获取文章", http.StatusInternalServerError)
		return
	}

	// 获取当前用户
	user := utils.GetUserFromSession(r)

	// 渲染模板
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/posts/list.html",
	)
	if err != nil {
		http.Error(w, "模板解析错误", http.StatusInternalServerError)
		return
	}

	// 获取当前年份
	currentYear := time.Now().Year()

	data := map[string]interface{}{
		"Title":       "所有文章",
		"Posts":       posts,
		"User":        user,
		"CurrentYear": currentYear,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "模板渲染错误", http.StatusInternalServerError)
	}
}

// GetPostHandler 处理单个文章请求
func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	// 从URL中提取文章ID
	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 获取存储实例
	store, err := db.NewSQLiteStore("./goblog.db")
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	defer store.Close()

	// 获取文章
	post, err := store.FindPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 获取当前用户
	user := utils.GetUserFromSession(r)

	// 渲染模板
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/posts/show.html",
	)
	if err != nil {
		http.Error(w, "模板解析错误", http.StatusInternalServerError)
		return
	}

	// 获取当前年份
	currentYear := time.Now().Year()

	data := map[string]interface{}{
		"Title":       post.Title,
		"Post":        post,
		"User":        user,
		"CurrentYear": currentYear,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "模板渲染错误", http.StatusInternalServerError)
	}
}

// NewPostFormHandler 处理新文章表单请求
func NewPostFormHandler(w http.ResponseWriter, r *http.Request) {
	// 检查用户是否已登录
	user := utils.GetUserFromSession(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 渲染模板
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/posts/new.html",
	)
	if err != nil {
		http.Error(w, "模板解析错误", http.StatusInternalServerError)
		return
	}

	// 获取当前年份
	currentYear := time.Now().Year()

	data := map[string]interface{}{
		"Title":       "创建新文章",
		"User":        user,
		"CurrentYear": currentYear,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "模板渲染错误", http.StatusInternalServerError)
	}
}

// CreatePostHandler 处理创建文章请求
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// 检查用户是否已登录
	user := utils.GetUserFromSession(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "表单解析错误", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	title := r.FormValue("title")
	content := r.FormValue("content")

	// 简单验证
	if title == "" || content == "" {
		http.Error(w, "标题和内容不能为空", http.StatusBadRequest)
		return
	}

	// 创建文章
	post := &models.Post{
		Title:   title,
		Content: content,
		UserID:  user.ID,
	}

	// 获取存储实例
	store, err := db.NewSQLiteStore("./goblog.db")
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	defer store.Close()

	// 保存文章
	if err := store.CreatePost(post); err != nil {
		http.Error(w, "无法创建文章", http.StatusInternalServerError)
		return
	}

	// 重定向到文章页面
	http.Redirect(w, r, "/posts/"+strconv.Itoa(post.ID), http.StatusSeeOther)
}

// EditPostFormHandler 处理编辑文章表单请求
func EditPostFormHandler(w http.ResponseWriter, r *http.Request) {
	// 检查用户是否已登录
	user := utils.GetUserFromSession(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 从URL中提取文章ID
	path := strings.TrimPrefix(r.URL.Path, "/posts/edit/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 获取存储实例
	store, err := db.NewSQLiteStore("./goblog.db")
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	defer store.Close()

	// 获取文章
	post, err := store.FindPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 检查是否是文章作者
	if post.UserID != user.ID {
		http.Error(w, "没有权限编辑该文章", http.StatusForbidden)
		return
	}

	// 渲染模板
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/posts/edit.html",
	)
	if err != nil {
		http.Error(w, "模板解析错误", http.StatusInternalServerError)
		return
	}

	// 获取当前年份
	currentYear := time.Now().Year()

	data := map[string]interface{}{
		"Title":       "编辑文章",
		"Post":        post,
		"User":        user,
		"CurrentYear": currentYear,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "模板渲染错误", http.StatusInternalServerError)
	}
}

// UpdatePostHandler 处理更新文章请求
func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	// 检查用户是否已登录
	user := utils.GetUserFromSession(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 从URL中提取文章ID
	path := strings.TrimPrefix(r.URL.Path, "/posts/update/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 获取存储实例
	store, err := db.NewSQLiteStore("./goblog.db")
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	defer store.Close()

	// 获取文章
	post, err := store.FindPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 检查是否是文章作者
	if post.UserID != user.ID {
		http.Error(w, "没有权限编辑该文章", http.StatusForbidden)
		return
	}

	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "表单解析错误", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	title := r.FormValue("title")
	content := r.FormValue("content")

	// 简单验证
	if title == "" || content == "" {
		http.Error(w, "标题和内容不能为空", http.StatusBadRequest)
		return
	}

	// 更新文章数据
	post.Title = title
	post.Content = content

	// 保存文章
	if err := store.UpdatePost(post); err != nil {
		http.Error(w, "无法更新文章", http.StatusInternalServerError)
		return
	}

	// 重定向到文章页面
	http.Redirect(w, r, "/posts/"+strconv.Itoa(post.ID), http.StatusSeeOther)
}

// DeletePostHandler 处理删除文章请求
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// 检查用户是否已登录
	user := utils.GetUserFromSession(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 从URL中提取文章ID
	path := strings.TrimPrefix(r.URL.Path, "/posts/delete/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 获取存储实例
	store, err := db.NewSQLiteStore("./goblog.db")
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	defer store.Close()

	// 获取文章
	post, err := store.FindPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 检查是否是文章作者
	if post.UserID != user.ID {
		http.Error(w, "没有权限删除该文章", http.StatusForbidden)
		return
	}

	// 删除文章
	if err := store.DeletePost(id); err != nil {
		http.Error(w, "无法删除文章", http.StatusInternalServerError)
		return
	}

	// 重定向到文章列表
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
