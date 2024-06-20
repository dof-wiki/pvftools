package model

type Dungeon struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func (d *Dungeon) GetCode() int {
	return d.ID
}

func (d *Dungeon) GetName() string {
	return d.Name
}
