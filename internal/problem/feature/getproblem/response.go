package getproblem

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

type ResponseProblemDetail struct {
	Language     string `json:"language"`
	Title        string `json:"title"`
	Background   string `json:"background"`
	Statement    string `json:"statement"`
	InputFormat  string `json:"input_format"`
	OutputFormat string `json:"output_format"`
	Note         string `json:"note"`
}

type ResponseProblemExample struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type ResponseProblemVersion struct {
	VersionID         uuid.UUID                `json:"version_id"`
	ProblemDifficulty dto.ProblemDifficulty    `json:"problem_difficulty"`
	Details           []ResponseProblemDetail  `json:"details"`
	Examples          []ResponseProblemExample `json:"examples"`
	CreatedAt         time.Time                `json:"created_at"`
}

type ResponseProblem struct {
	ProblemID       uuid.UUID                `json:"problem_id"`
	LatestVersionID uuid.UUID                `json:"latest_version_id"`
	Versions        []ResponseProblemVersion `json:"versions"`
	Status          constant.ProblemStatus   `json:"status"`
	Creator         ResponseUser             `json:"creator"`
	Reviewer        *ResponseUser            `json:"reviewer"`
	Testers         []ResponseUser           `json:"testers"`
	TargetContest   *ResponseContest         `json:"target_contest"`
	AssignedContest *ResponseContest         `json:"assigned_contest"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

type Response struct {
	Problem ResponseProblem `json:"problem"`
}
