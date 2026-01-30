package model

type Comments struct {
	BaseModel
	Content string `gorm:"type:varchar(100)"`
	UserId  uint   `gorm:"index"`
	PostId  uint   `gorm:"index"`
}
