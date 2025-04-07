package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recover 恢复中间件 - 处理程序崩溃
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误和堆栈信息
				log.Printf("程序崩溃: %v\n", err)
				log.Printf("堆栈信息: %s\n", debug.Stack())

				// 返回500内部服务器错误
				http.Error(w, "内部服务器错误", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
