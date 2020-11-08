package model

import (
	"time"

	"github.com/JingruiLea/gobatis/example/common"

	"gorm.io/gorm"
)

type User struct {
	ID        int64 `gorm:"autoIncrement;primary_key"`
	UserID    int64
	Username  string
	Age       int32
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}

type BO struct {
}

func (User) TableName() string {
	return "user"
}

func (u *User) FromBO(BO *common.User) {
	*u = User{}
}
func (u *User) ToBO() *common.User {
	return &common.User{}
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
	return tx.Create(&User{
		ID:        0,
		UserID:    0,
		Username:  "小明",
		Age:       0,
		Phone:     "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		IsDeleted: false,
	}).Error
}
