package model

type CustomAttrTmpl struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Tmpl    string `json:"tmpl"`
	Desc    string `json:"desc"`
	Choices string `json:"choices"`
}
