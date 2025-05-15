package constant

type ProblemStatus string

const (
	ProblemStatusRejected           ProblemStatus = "rejected"
	ProblemStatusPendingReview      ProblemStatus = "pending_review"
	ProblemStatusNeedsRevision      ProblemStatus = "needs_revision"
	ProblemStatusApprovedForTesting ProblemStatus = "approved_for_testing"
	ProblemStatusAwaitingFinalCheck ProblemStatus = "awaiting_final_check"
)

func FromStringToProblemStatus(status string) ProblemStatus {
	switch status {
	case string(ProblemStatusRejected):
		return ProblemStatusRejected
	case string(ProblemStatusPendingReview):
		return ProblemStatusPendingReview
	case string(ProblemStatusNeedsRevision):
		return ProblemStatusNeedsRevision
	case string(ProblemStatusApprovedForTesting):
		return ProblemStatusApprovedForTesting
	case string(ProblemStatusAwaitingFinalCheck):
		return ProblemStatusAwaitingFinalCheck
	default:
		return ""
	}
}
