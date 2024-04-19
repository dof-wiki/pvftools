package model

type Stackable struct {
	Code int    `json:"code" gorm:"primaryKey"`
	Name string `json:"name" gorm:"index"`
	Path string `json:"path" gorm:"index"`
}
