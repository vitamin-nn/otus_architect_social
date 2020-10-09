package error

type OutError string

func (e OutError) Error() string {
	return string(e)
}

var (
	ErrUserAlreadyExists   = OutError("user with this email already registered")
	ErrFriendAlreadyExists = OutError("this friend was already added")
)
