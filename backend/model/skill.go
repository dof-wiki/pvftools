package model

type Skill struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Job  int    `json:"job"`
	Code int    `json:"code" gorm:"index"`
	Name string `json:"name" gorm:"index"`
	Path string `json:"path"`
}

func (s *Skill) GetCode() int {
	return s.Code
}

func (s *Skill) GetName() string {
	return s.Name
}
