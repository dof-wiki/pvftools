package proto

type SkillIncrement struct {
	Job           string `json:"job"`
	SkillId       int    `json:"skill_id"`
	Type          string `json:"type"`
	DataType      string `json:"data_type"`
	DataIndex     int    `json:"data_index"`
	IncrementType string `json:"increment_type"`
	Value         int    `json:"value"`
}

func (s *SkillIncrement) FromData(data [7]any) {
	s.Job = data[0].(string)
	s.SkillId = data[1].(int)
	s.Type = data[2].(string)
	s.DataType = data[3].(string)
	s.DataIndex = data[4].(int)
	s.IncrementType = data[5].(string)
	s.Value = data[6].(int)
}
