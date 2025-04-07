package controllers

import (
	"goblog/db"
	"goblog/models"
	"goblog/utils"
	"html/template"
	"log"
	"net/http"
	"time"
)

// HomeHandler 处理首页请求
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	log.Println("处理首页请求")

	// 获取当前用户
	user := utils.GetUserFromSession(r)
	log.Printf("当前用户: %v", user)

	// 预设空文章列表
	posts := []*models.Post{}

	// 尝试获取文章，但如果失败也继续显示页面
	store, err := db.NewSQLiteStore("./goblog.db")
	if err != nil {
		log.Printf("数据库连接错误: %v", err)
	} else {
		defer store.Close()

		// 尝试获取文章
		retrievedPosts, err := store.FindAllPosts()
		if err != nil {
			log.Printf("获取文章错误: %v", err)
		} else if retrievedPosts != nil {
			posts = retrievedPosts
		}
	}

	log.Printf("获取到 %d 篇文章", len(posts))

	// 渲染模板 - 修复模板解析方式
	files := []string{
		"templates/base.html",
		"templates/home.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("模板解析错误: %v", err)
		http.Error(w, "模板解析错误", http.StatusInternalServerError)
		return
	}

	// 获取当前年份
	currentYear := time.Now().Year()

	data := map[string]interface{}{
		"Title":       "博客首页",
		"Posts":       posts,
		"User":        user,
		"CurrentYear": currentYear,
	}

	log.Println("开始渲染模板")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("模板渲染错误: %v", err)
		http.Error(w, "模板渲染错误", http.StatusInternalServerError)
	}
	log.Println("首页请求处理完成")
}
