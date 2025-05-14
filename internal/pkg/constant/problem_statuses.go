package constant

type ProblemStatus string

const (
	ProblemStatusPendingReview ProblemStatus = "pending_review"
	ProblemStatusNeedsRevision ProblemStatus = "needs_revision"
)

func FromStringToProblemStatus(status string) ProblemStatus {
	switch status {
	case string(ProblemStatusPendingReview):
		return ProblemStatusPendingReview
	case string(ProblemStatusNeedsRevision):
		return ProblemStatusNeedsRevision
	default:
		return ""
	}
}
