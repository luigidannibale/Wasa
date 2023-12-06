package database

type ConstError string

func (err ConstError) Error() string {
	return string(err)
}

const (
	ErrNotFound            = ConstError("NotFound")
	ErrInternalServerError = ConstError("InternalServerError")
	ErrAlreadyDone         = ConstError("AlreadyDone")
	ErrUsernameTaken       = ConstError("UsernameTaken")
)
