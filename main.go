package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

var (
	DB *gorm.DB
)

func initMySQL() (err error) {
	dsn := "root:20010111@tcp(127.0.0.1:3306)/dbtest?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置连接的最大存活时间
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	return nil
}

func main() {
	// 创建数据库
	err := initMySQL()
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&Todo{})

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
		v1Group.POST("/todo", func(ctx *gin.Context) {
			var todo Todo
			ctx.BindJSON(&todo)
			err := DB.Create(&todo).Error
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				ctx.JSON(http.StatusOK, todo)
			}
		})
		// 查看
		v1Group.GET("/todo", func(ctx *gin.Context) {
			// 查询表中所有的数据
			var todolist []Todo
			err := DB.Find(&todolist).Error
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				ctx.JSON(http.StatusOK, todolist)
			}

		})
		v1Group.GET("/todo/:id", func(ctx *gin.Context) {

		})
		// 修改
		v1Group.PUT("/todo/:id", func(ctx *gin.Context) {
			id, ok := ctx.Params.Get("id")
			if !ok {
				ctx.JSON(http.StatusOK, gin.H{
					"error": "error!",
				})
				return
			}
			var todo Todo
			if err := DB.Where("id=?", id).First(&todo).Error; err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			}
			ctx.BindJSON(&todo)
			if err = DB.Save(&todo).Error; err != nil {
				ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, todo)
			}

		})
		// 删除
		v1Group.DELETE("/todo/:id", func(ctx *gin.Context) {
			id, ok := ctx.Params.Get("id")
			if !ok {
				ctx.JSON(http.StatusOK, gin.H{
					"error": "error!",
				})
				return
			}
			if err := DB.Where("id=?", id).Delete(Todo{}).Error; err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{id: "delete"})
			}
		})
	}

	r.Run(":9000")
}
