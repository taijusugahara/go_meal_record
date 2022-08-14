package model

import (
	"encoding/json"
	"gorm.io/gorm"
)

type User struct {
	//gorm.Modelを付与するとcreated?at,updated_at,deleted_at自動取得
	// 個別にcreated_atだけ取得したい場合はCreatedAT time.Timeと書く
	gorm.Model `json:"-"` //json非表示(created_at,updated_at,deleted_at)
	ID         int        `gorm:"primary_key" json:"ID"`
	Name       string     `validate:"required,max=255" json:"name"`
	Email      string     `gorm:"unique" validate:"required,email,max=255" json:"email,omitempty"`
	Password   string     `validate:"required" json:"password,omitempty"`
}

// //json返す時passwordを空に 上のjson:"-"ではpostの時値渡せない
func (u User) MarshalJSON() ([]byte, error) {
	type user User // prevent recursion
	x := user(u)
	x.Email = ""
	x.Password = ""
	return json.Marshal(x)
}
