package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func main() {
	// 创建数据库

	// 连接数据库

	r := gin.Default()
	// 寻找静态文件，当找 /static,就去static目录下找文件
	r.Static("/static", "static")
	// 寻找html文件
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	v1Group := r.Group("v1")
	{
		// 代办事项

		// 添加
		v1Group.GET("/todo", func(ctx *gin.Context) {

		})
		// 查看
		v1Group.GET("/todo", func(ctx *gin.Context) {

		})
		v1Group.GET("/todo/:id", func(ctx *gin.Context) {

		})
		// 修改
		v1Group.POST("/todo/:id", func(ctx *gin.Context) {

		})
		// 删除
		v1Group.DELETE("/todo/:id", func(ctx *gin.Context) {

		})
	}

	r.Run(":9000")
}
