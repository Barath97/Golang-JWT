package models

import "time"

type User struct {
	ID            uint      `gorm:"primarykey;autoIncrement"  json:"id" validate:"required"`
	First_name    *string   `gorm:"size:200;not null" json:"first_name" validate:"required"`
	Last_name     *string   `gorm:"size:200" json:"last_name" validate:"required"`
	Password      *string   `gorm:"size:50;not null" json:"-" validate:"required"`
	Email         *string   `gorm:"size:100;unique;not null" json:"email" validate:"required,email"`
	Phone         *string   `gorm:"size:15;unique" json:"phone" validate:"required"`
	Token         *string   `json:"token"`
	User_type     *string   `gorm:"size:50;default:'user'" json:"user_type"`
	Refresh_token *string   `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	User_id       *string   `gorm:"unique;not null" json:"user_id"`
}
