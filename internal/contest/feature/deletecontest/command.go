package deletecontest

type Command struct {
	ContestID string `param:"contest_id" validate:"required"`
}
