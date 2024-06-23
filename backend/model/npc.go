package model

type Npc struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func (n *Npc) GetCode() int {
	return n.ID
}

func (n *Npc) GetName() string {
	return n.Name
}
