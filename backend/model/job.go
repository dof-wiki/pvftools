package model

type Job struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Code         int    `json:"code"`
	Job          string `json:"job"`
	GrowType     int    `json:"grow_type"`
	GrowTypeName string `json:"grow_type_name"`

	JobName string `json:"job_name" gorm:"-"`
}
