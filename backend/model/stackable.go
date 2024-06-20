package model

type Stackable struct {
	Code int    `json:"code" gorm:"primaryKey"`
	Name string `json:"name" gorm:"index"`
	Path string `json:"path" gorm:"index"`
}

func (s *Stackable) GetCode() int {
	return s.Code
}

func (s *Stackable) GetName() string {
	return s.Name
}
