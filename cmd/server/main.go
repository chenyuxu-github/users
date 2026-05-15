package main

import (
	"log"

	"user/internal/config"
	"user/internal/database"
	"user/internal/router"
)

func main() {
	// 读取项目配置，包括服务监听地址和数据库连接信息。
	cfg, err := config.Load("config/config.toml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 根据配置初始化数据库连接，后续 handler 统一复用这个连接对象。
	db, err := database.Init(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("init database: %v", err)
	}

	// 注册静态资源和 API 路由，并按配置的地址启动 HTTP 服务。
	r := router.Setup(db)
	if err := r.Run(cfg.Server.Address()); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
