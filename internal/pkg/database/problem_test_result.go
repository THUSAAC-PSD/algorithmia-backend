package database

import (
	"time"

	"github.com/google/uuid"
)

type ProblemTestResult struct {
	ProblemTestResultID uuid.UUID `gorm:"primaryKey;type:uuid"`
	VersionID           uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_problem_test_result_version_tester"`
	TesterID            uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_problem_test_result_version_tester"`
	Tester              User      `gorm:"foreignKey:TesterID"`
	Status              string
	Comment             string
	CreatedAt           time.Time
}
