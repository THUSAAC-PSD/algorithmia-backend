package listproblem

import (
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/dto"

	"github.com/google/uuid"
)

type ResponseProblemTitle struct {
	Language string `json:"language"`
	Title    string `json:"title"`
}

type ResponseUser struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}

type ResponseContest struct {
	ContestID uuid.UUID `json:"contest_id"`
	Title     string    `json:"title"`
}

type ResponseProblem struct {
	ProblemID         uuid.UUID              `json:"problem_id"`
	Titles            []ResponseProblemTitle `json:"title"`
	Status            constant.ProblemStatus `json:"status"`
	Creator           ResponseUser           `json:"creator"`
	Reviewer          *ResponseUser          `json:"reviewer"`
	Testers           []ResponseUser         `json:"testers"`
	TargetContest     *ResponseContest       `json:"target_contest"`
	AssignedContest   *ResponseContest       `json:"assigned_contest"`
	ProblemDifficulty dto.ProblemDifficulty  `json:"problem_difficulty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type Response struct {
	Problems []ResponseProblem `json:"problems"`
}
