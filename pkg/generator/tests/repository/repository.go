package repository

type User struct {
	ID   int
	Name string
}

//go:generate ffakes -i UserRepository
type UserRepository interface {
	GetAllUsers() ([]User, error)
	FindUserByID(id int) (u User, err error)
	CreateUser(user User) error
	DeleteUserByID(id int) error
	UpdateUser(old, new User) error
	Execute(user User) error
	FindUserByName(name string) (User, error)
}
