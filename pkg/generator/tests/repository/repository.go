package repository

import "context"

type User struct {
	ID   int
	Name string
}

//go:generate ffakes -i UserRepository
type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]User, error)
	FindUserByID(id int) (u User, err error)
	CreateUser(user User) error
	DeleteUserByID(id int) error
	UpdateUser(ctx context.Context, old, new User) error
	Execute(user User) error
	FindUserByName(name string) (User, error)
}
