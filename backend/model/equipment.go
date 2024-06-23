package model

type Equipment struct {
	Code      int    `json:"code" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"index"`
	Path      string `json:"path" gorm:"index"`
	Rarity    int    `json:"rarity"`
	Type      string `json:"type"`
	PartSet   int    `json:"part_set"`
	MiniLevel int    `json:"mini_level"`
}

func (e *Equipment) GetCode() int {
	return e.Code
}

func (e *Equipment) GetName() string {
	return e.Name
}
