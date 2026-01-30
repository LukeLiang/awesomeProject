package initialize

import (
	"awesomeProject/homework04/common"
	"awesomeProject/homework04/global"
	"awesomeProject/homework04/model"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// LoadSystemConfig 加载配置文件
func LoadSystemConfig() {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")            // 当前目录
	v.AddConfigPath("./homework04") // homework04 目录

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic("加载配置文件失败")
	}

	// 解析配置到结构体
	var config common.Config
	if err := v.Unmarshal(&config); err != nil {
		panic("解析配置文件失败")
	}

	// 设置 JWT Secret（转换为 []byte）
	config.Jwt.SetSecret(v.GetString("jwt.secret"))

	global.Config = &config
}

// OpenDatabase 开启数据库连接
func OpenDatabase() {

	if global.Config.Type != string(common.MySQL) {
		panic("不支持的数据库类型")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.Config.User,
		global.Config.Password,
		global.Config.Host,
		global.Config.DB.Port,
		global.Config.Database)

	// 配置日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true, // 禁用物理外键
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 设置为 true，则 User 对应表名为 user
		},
		Logger: newLogger,
	})
	if err != nil {
		global.Logger.Error("failed to connect database:", zap.Error(err))
	}

	fmt.Println("数据库连接成功....")

	global.DB = db
}

// InitDatabase 初始化数据库表
func InitDatabase() {
	users := []model.User{
		//{Username: "admin", Password: "admin", Email: "admin@example.com"},
		//{Username: "user1", Password: "user1", Email: "user1@example.com"},
		//{Username: "user2", Password: "user2", Email: "user2@example.com"},
	}
	posts := []model.Posts{
		//{Title: "Posts 1", Content: "This is the first post.", UserId: 1},
		//{Title: "Posts 2", Content: "This is the second post.", UserId: 1},
		//{Title: "Posts 3", Content: "This is the third post.", UserId: 2},
	}
	comments := []model.Comments{
		//{Content: "Comments 1", UserId: 1, PostId: 1},
		//{Content: "Comments 2", UserId: 1, PostId: 1},
		//{Content: "Comments 3", UserId: 2, PostId: 2},
	}

	_ = global.DB.AutoMigrate(&users, &posts, &comments)
}
