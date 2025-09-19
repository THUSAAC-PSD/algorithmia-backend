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

	PermissionProblemDraftCreate                  = "problem:draft:create"
	PermissionProblemDraftReadOwn                 = "problem:draft:read:own"
	PermissionProblemDraftUpdateOwn               = "problem:draft:update:own"
	PermissionProblemDraftDeleteOwn               = "problem:draft:delete:own"
	PermissionProblemDraftSubmitOwn               = "problem:draft:submit:own"
	PermissionProblemListAll                      = "problem:list:all"
	PermissionProblemListCreatedOwn               = "problem:list:created-own"
	PermissionProblemListAwaitingReviewAll        = "problem:list:awaiting-review-all"
	PermissionProblemListAssignedTest             = "problem:list:assigned-test"
	PermissionProblemReadDetailsAny               = "problem:read:details:any"
	PermissionProblemReadDetailsCreatedOwn        = "problem:read:details:created-own"
	PermissionProblemReadDetailsAwaitingReviewAny = "problem:read:details:awaiting-review-any"
	PermissionProblemReadDetailsAssignedTest      = "problem:read:details:assigned-test"
	PermissionProblemReviewAny                    = "problem:review:any"
	PermissionProblemReviewOverride               = "problem:review:override"
	PermissionProblemAssignTesters                = "problem:assign:testers"
	PermissionProblemTestAssigned                 = "problem:test:assigned"
	PermissionProblemTestOverride                 = "problem:test:override"

	PermissionContestListAll            = "contest:list_all"
	PermissionContestReadDetailsAny     = "contest:read_details_any"
	PermissionContestCreate             = "contest:create"
	PermissionContestUpdateAny          = "contest:update_any"
	PermissionContestDeleteAny          = "contest:delete_any"
	PermissionContestAssignProblemAny   = "contest:assign_problem_any"
	PermissionContestUnassignProblemAny = "contest:unassign_problem_any"
)
