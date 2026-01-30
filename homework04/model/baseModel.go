package model

import (
	"database/sql"

	"awesomeProject/homework04/common/types"
)

type BaseModel struct {
	ID        uint            `gorm:"primarykey;autoincrement" json:"id"`
	CreatedAt types.LocalTime `json:"createdAt"`
	UpdatedAt types.LocalTime `json:"updatedAt"`
	DeletedAt sql.NullTime    `gorm:"index" json:"-"`
}
