package database

type ConstError string

func (err ConstError) Error() string {
	return string(err)
}

const (
	NotFound            = ConstError("NotFound")
	InternalServerError = ConstError("InternalServerError")
	AlreadyDone         = ConstError("AlreadyDone")
	UsernameTaken       = ConstError("UsernameTaken")
)
