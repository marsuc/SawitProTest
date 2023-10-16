// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)
	CreateUser(ctx context.Context, user User) (id int, err error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user User, err error)
	GetUserById(ctx context.Context, id int) (user User, err error)
	UpdateUser(ctx context.Context, user User) (err error)
}
