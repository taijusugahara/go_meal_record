package model

import (
	"encoding/json"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int    `gorm:"primary_key" json:"ID"`
	Name     string `validate:"required,max=255" json:"name"`
	Email    string `gorm:"unique" validate:"required,email,max=255" json:"email"`
	Password string `validate:"required" json:"password,omitempty"`
}

//json返す時passwordを空に
func (u User) MarshalJSON() ([]byte, error) {
	type user User // prevent recursion
	x := user(u)
	x.Password = ""
	return json.Marshal(x)
}
