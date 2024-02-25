package entities

type StorageError struct {
	Err error
}

func (e StorageError) Error() string {
	return e.Err.Error()
}

type NotFoundError struct{}

func (e NotFoundError) Error() string {
	return "not found"
}

type InvalidPasswordError struct{}

func (e InvalidPasswordError) Error() string {
	return "invalid password"
}

type InvalidUserError struct{}

func (e InvalidUserError) Error() string {
	return "invalid user"
}

type ExistUserError struct{}

func (e ExistUserError) Error() string {
	return "exist user"
}

type NotExistUserError struct{}

func (e NotExistUserError) Error() string {
	return "not exist user"
}

type UnexpectedError struct {
	Err error
}

func (e UnexpectedError) Error() string {
	return e.Err.Error()
}

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	if e.Err == nil {
		return "validation error"
	}
	return e.Err.Error()
}

type InvalidOrderError struct{}

func (e InvalidOrderError) Error() string {
	return "invalid order"
}

type ExistOrderError struct{}

func (e ExistOrderError) Error() string {
	return "exist order"
}

type OrderIsCreatedByAnotherUserError struct{}

func (e OrderIsCreatedByAnotherUserError) Error() string {
	return "order is created by another user"
}

type OutOfBalanceError struct{}

func (e OutOfBalanceError) Error() string {
	return "out of balance"
}
