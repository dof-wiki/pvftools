package box

type BoxItem struct {
	Id    int `json:"id"`
	Count int `json:"count"`
	Rate  int `json:"rate"`
}

type BoxItemGroup struct {
	Items []*BoxItem `json:"items"`
}

const (
	BoxTypeRandom = iota
	BoxTypeSelection
)

type BoxParams struct {
	BoxType       int             `json:"box_type"`
	Name          string          `json:"name"`
	Explain       string          `json:"explain"`
	Grade         int             `json:"grade"`
	Rarity        int             `json:"rarity"`
	Items         []*BoxItemGroup `json:"items"`
	ItemType      string          `json:"tp"` // 物品类型
	CategoryNames []string        `json:"category_names"`
}
