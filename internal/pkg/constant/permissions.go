package constant

const (
	PermissionMediaUploadForDraftOwn = "media:upload_for_draft_own"
	PermissionMediaUploadForChatOwn  = "media:upload_for_chat_own"

	PermissionUserListAll          = "user:list_all"
	PermissionUserReadProfileAny   = "user:read_profile_any"
	PermissionUserReadProfileOwn   = "user:read_profile_own"
	PermissionUserUpdateProfileOwn = "user:update_profile_own"
	PermissionUserManageRolesAny   = "user:manage_roles_any"

	PermissionRoleListAll   = "role:list_all"
	PermissionRoleManageAny = "role:manage_any"

	PermissionProblemDraftCreate                  = "problem_draft:create"
	PermissionProblemDraftReadOwn                 = "problem_draft:read_own"
	PermissionProblemDraftUpdateOwn               = "problem_draft:update_own"
	PermissionProblemDraftDeleteOwn               = "problem_draft:delete_own"
	PermissionProblemDraftSubmitOwn               = "problem_draft:submit_own"
	PermissionProblemListAll                      = "problem:list_all"
	PermissionProblemListCreatedOwn               = "problem:list_created_own"
	PermissionProblemListAwaitingReviewAll        = "problem:list_awaiting_review_all"
	PermissionProblemListAssignedTest             = "problem:list_assigned_test"
	PermissionProblemReadDetailsAny               = "problem:read_details_any"
	PermissionProblemReadDetailsCreatedOwn        = "problem:read_details_created_own"
	PermissionProblemReadDetailsAwaitingReviewAny = "problem:read_details_awaiting_review_any"
	PermissionProblemReadDetailsAssignedTest      = "problem:read_details_assigned_test"
	PermissionProblemReviewAny                    = "problem:review_any"
	PermissionProblemReviewOverride               = "problem:review_override"
	PermissionProblemAssignTesters                = "problem:assign_testers"
	PermissionProblemTestAssigned                 = "problem:test_assigned"
	PermissionProblemTestOverride                 = "problem:test_override"

	PermissionContestListAll            = "contest:list_all"
	PermissionContestReadDetailsAny     = "contest:read_details_any"
	PermissionContestCreate             = "contest:create"
	PermissionContestUpdateAny          = "contest:update_any"
	PermissionContestDeleteAny          = "contest:delete_any"
	PermissionContestAssignProblemAny   = "contest:assign_problem_any"
	PermissionContestUnassignProblemAny = "contest:unassign_problem_any"
)
