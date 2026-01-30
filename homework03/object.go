package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type User struct {
	gorm.Model
	Name     string `gorm:"unique;type:varchar(100)"`
	Age      int
	Gender   int    `gorm:"type:tinyint"`
	Posts    []Post `gorm:"foreignKey:UserID"`
	PostNums uint   `gorm:"type:int unsigned;default:0;not null"`
}

type Post struct {
	gorm.Model
	Title    string `gorm:"type:varchar(100)"`
	Author   string `gorm:"type:varchar(100)"`
	Like     uint
	UserID   uint      `gorm:"index"`
	Comments []Comment `gorm:"foreignKey:PostId"`
	Status   string    `gorm:"type:varchar(50)"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:varchar(100)"`
	Like    uint
	PostId  uint `gorm:"index"`
}

// AfterDelete 在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"
func (comment *Comment) AfterDelete(db *gorm.DB) (err error) {
	var count int64
	err = db.Model(&Comment{}).Where("post_id", comment.PostId).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		db.Model(&Post{}).Where("id", comment.PostId).Update("status", "无评论")
	}

	return nil
}

func (post *Post) AfterCreate(db *gorm.DB) (err error) {
	user := User{}
	db.Model(&user).
		Where("id", post.UserID).
		Update("post_nums", gorm.Expr("post_nums + ?", 1))

	return nil
}

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

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
		panic("failed to connect database: " + err.Error())
	}

	fmt.Printf("数据库连接成功....")

	//user := User{}
	//post := Post{}
	//comment := Comment{}
	//_ = db.AutoMigrate(&user, &post, &comment)
	//
	//users := []User{
	//	{Name: "Luke1", Age: 10, Gender: 1},
	//	{Name: "Luke2", Age: 27, Gender: 1},
	//	{Name: "cherry", Age: 18, Gender: 2},
	//}
	//
	//posts := []Post{
	//	{Title: "this is a title1", Author: "Luke1", Like: 0, UserID: 1, Status: "无状态"},
	//	{Title: "this is a title2", Author: "Luke1", Like: 10, UserID: 1, Status: "无状态"},
	//	{Title: "this is a title3", Author: "Luke1", Like: 0, UserID: 1, Status: "无状态"},
	//}
	//
	//comments := []Comment{
	//	{Content: "hello1", Like: 0, PostId: 2},
	//	{Content: "hello2", Like: 0, PostId: 2},
	//	{Content: "hello3", Like: 0, PostId: 2},
	//}
	//
	//db.Create(&users)
	//db.Create(&posts)
	//db.Create(&comments)

	// 查询某个用户发布的所有文章及其对应的评论信息。
	//var post []Post
	//err = db.Model(&Post{}).Preload("Comments", &[]Comment{}).Find(&post, "user_id", 1).Error
	//if err != nil {
	//	fmt.Printf(err.Error())
	//}

	// 查询评论数量最多的文章信息。
	//maps := map[string]interface{}{}
	//err = db.Model(&Comment{}). // 聚合查询通常需要显式指定表名
	//				Select("post_id, count(id) as total").
	//				Group("post_id").
	//				Order("total DESC"). // 使用别名排序
	//				Take(&maps).Error
	//
	//postId := maps["post_id"]
	//post := Post{}
	//err = db.Preload("Comments").First(&post, postId).Error
	//
	//err = db.First(&post, postId).Error

	// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	//post := Post{
	//	Title:  "this is hook test for post",
	//	Author: "Luke",
	//	Like:   0,
	//	UserID: 2,
	//	Status: "无状态",
	//}
	//db.Create(&post)

	// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

	comment := Comment{
		Model:  gorm.Model{ID: 4},
		PostId: 4,
	}
	db.Delete(&comment)
}
