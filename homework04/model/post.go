package model

type Posts struct {
	BaseModel
	Title    string     `gorm:"type:varchar(100)" validate:"required"`
	Content  string     `gorm:"type:varchar(255)" validate:"required"`
	UserId   uint       `gorm:"index"`
	Comments []Comments `gorm:"foreignKey:PostId"`
}
