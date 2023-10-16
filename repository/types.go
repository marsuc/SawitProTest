// This file contains types that are used in the repository layer.
package repository

import (
	"time"

	"github.com/marsuc/SawitProTest/generated"
)

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type User struct {
	Id          int
	PhoneNumber string `validate:"required,phoneNumber"`
	FullName    string `validate:"required,min=3,max=60"`
	Password    string `validate:"required,min=6,max=64,password"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLogin   *time.Time
	LoginCount  int
}

func (u User) FromRegisterRequest(req generated.RegisterRequest) User {
	if req.PhoneNumber != nil {
		u.PhoneNumber = *req.PhoneNumber
	}

	if req.FullName != nil {
		u.FullName = *req.FullName
	}

	if req.Password != nil {
		u.Password = *req.Password
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return u
}
