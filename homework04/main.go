package main

import (
	"awesomeProject/homework04/initialize"
	"awesomeProject/homework04/router"
)

// 启动应用程序
func main() {
	// 1. 加载配置文件
	initialize.LoadSystemConfig()

	// 2. 初始化日志
	initialize.InitLogger()

	// 3. 连接数据库
	initialize.OpenDatabase()

	// 4. 创建数据库表
	initialize.InitDatabase()

	// 5. 启动服务
	router.InitServer()

}
