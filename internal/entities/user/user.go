package user

import (
	"errors"
)

type User struct {
	ID       int
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Status   Status
}
type Status string

// input string must be "active" or "offline"
func BuildStatus(status string) (Status, error) {
	if status == "active" {
		return Status(status), nil
	} else if status == "offline" {
		return Status(status), nil
	} else {
		return Status(""), errors.New("incorrect status:" + status)
	}
}
