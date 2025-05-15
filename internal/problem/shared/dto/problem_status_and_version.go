package dto

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"

	"github.com/google/uuid"
)

type ProblemStatusAndVersion struct {
	Status  constant.ProblemStatus
	DraftID uuid.UUID
}
