package storage

import "fmt"

type StorErr struct {
	At  string
	Err error
}

func (se StorErr) Error() string {
	return fmt.Sprintf("at %s: %v", se.At, se.Err)
}
func (se StorErr) Unwrap() error {
	return se.Err
}

type ErrSuchUserExists struct {
	StorErr
	Email string
}

func (e ErrSuchUserExists) Error() string {
	return fmt.Sprintf("user with %s already exists", e.Email)
}
func (e ErrSuchUserExists) Is(target error) bool {
	_, ok := target.(ErrSuchUserExists)
	return ok
}
