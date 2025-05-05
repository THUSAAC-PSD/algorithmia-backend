package database

import "github.com/google/uuid"

type ProblemDraftExample struct {
	ExampleID      uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProblemDraftID uuid.UUID `gorm:"type:uuid"`
	Input          string
	Output         string
}
