package router

import (
	"goblog/controllers"
	"goblog/middleware"
	"net/http"
)

// SetupRouter 设置路由
func SetupRouter() http.Handler {
	// 创建默认多路复用器
	mux := http.NewServeMux()

	// 静态文件服务
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// 首页
	mux.HandleFunc("/", controllers.HomeHandler)

	// 文章相关路由
	mux.HandleFunc("/posts", controllers.ListPostsHandler)
	mux.HandleFunc("/posts/", controllers.GetPostHandler)
	mux.HandleFunc("/posts/new", controllers.NewPostFormHandler)
	mux.HandleFunc("/posts/create", controllers.CreatePostHandler)
	mux.HandleFunc("/posts/edit/", controllers.EditPostFormHandler)
	mux.HandleFunc("/posts/update/", controllers.UpdatePostHandler)
	mux.HandleFunc("/posts/delete/", controllers.DeletePostHandler)

	// 用户相关路由
	mux.HandleFunc("/login", controllers.LoginFormHandler)
	mux.HandleFunc("/login/process", controllers.LoginProcessHandler)
	mux.HandleFunc("/logout", controllers.LogoutHandler)
	mux.HandleFunc("/register", controllers.RegisterFormHandler)
	mux.HandleFunc("/register/process", controllers.RegisterProcessHandler)

	// 应用中间件
	var handler http.Handler = mux
	handler = middleware.Logger(handler)
	handler = middleware.Recover(handler)

	return handler
}
