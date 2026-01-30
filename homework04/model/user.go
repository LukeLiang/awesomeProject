package model

type User struct {
	BaseModel
	Username string  `gorm:"type:varchar(100)"`
	Password string  `gorm:"type:varchar(100)"`
	Email    string  `gorm:"type:varchar(100)" validate:"email"`
	Posts    []Posts `gorm:"foreignKey:UserID"`
}
