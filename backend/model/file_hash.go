package model

type FileHash struct {
	Path string `json:"path" gorm:"primaryKey"`
	Md5  string `json:"md5"`
}
