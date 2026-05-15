package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"user/internal/handler"
)

// Setup 创建 Gin 引擎，注册前端静态资源和用户 API 路由。
func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 用户接口处理器依赖数据库连接，通过构造函数注入。
	userHandler := handler.NewUserHandler(db)

	// 托管前端页面和静态资源，访问 / 即可打开用户管理页面。
	r.StaticFile("/", "./web/index.html")
	r.Static("/assets", "./web/assets")

	// API 使用版本前缀，方便后续扩展 v2 或其他模块。
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", userHandler.GetUsers)
		v1.GET("/users/:id", userHandler.GetUserByID)
		v1.POST("/users", userHandler.CreateUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
	}

	return r
}
