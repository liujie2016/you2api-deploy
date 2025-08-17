package main

import (
	"fmt"
	"log"
	"net/http"

	api "you2api/api" // 请替换为您的实际项目名
)

func main() {
	// 将 /api/main.go 的 Handler 函数注册到根路径
	http.HandleFunc("/", api.Handler)

	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)

	// 启动服务器
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}
