package constant

type ProblemStatus string

const (
	ProblemStatusRejected                ProblemStatus = "rejected"
	ProblemStatusPendingReview           ProblemStatus = "pending_review"
	ProblemStatusNeedsRevision           ProblemStatus = "needs_revision"
	ProblemStatusPendingTesting          ProblemStatus = "pending_testing"
	ProblemStatusTestingChangesRequested ProblemStatus = "testing_changes_requested"
	ProblemStatusAwaitingFinalCheck      ProblemStatus = "awaiting_final_check"
	ProblemStatusCompleted               ProblemStatus = "completed"
)

func FromStringToProblemStatus(status string) ProblemStatus {
	switch status {
	case string(ProblemStatusRejected):
		return ProblemStatusRejected
	case string(ProblemStatusPendingReview):
		return ProblemStatusPendingReview
	case string(ProblemStatusNeedsRevision):
		return ProblemStatusNeedsRevision
	case string(ProblemStatusPendingTesting):
		return ProblemStatusPendingTesting
	case string(ProblemStatusTestingChangesRequested):
		return ProblemStatusTestingChangesRequested
	case string(ProblemStatusAwaitingFinalCheck):
		return ProblemStatusAwaitingFinalCheck
	case string(ProblemStatusCompleted):
		return ProblemStatusCompleted
	default:
		return ""
	}
}
