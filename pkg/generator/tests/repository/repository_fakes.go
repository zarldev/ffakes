// Code generated by ffakes v0.0.2 DO NOT EDIT.

package repository

import (
	"testing"
)

type FakeUserRepository struct {
	t                   *testing.T
	GetAllUsersCount    int
	FindUserByIDCount   int
	CreateUserCount     int
	DeleteUserByIDCount int
	UpdateUserCount     int
	ExecuteCount        int
	FindUserByNameCount int
	FGetAllUsers        []func() ([]User, error)
	FFindUserByID       []func(id int) (u User, err error)
	FCreateUser         []func(user User) error
	FDeleteUserByID     []func(id int) error
	FUpdateUser         []func(old, new User) error
	FExecute            []func(user User) error
	FFindUserByName     []func(name string) (User, error)
}

type GetAllUsersFunc = func() ([]User, error)
type FindUserByIDFunc = func(id int) (u User, err error)
type CreateUserFunc = func(user User) error
type DeleteUserByIDFunc = func(id int) error
type UpdateUserFunc = func(old, new User) error
type ExecuteFunc = func(user User) error
type FindUserByNameFunc = func(name string) (User, error)
type UserRepositoryOption func(f *FakeUserRepository)

func OnGetAllUsers(fn ...GetAllUsersFunc) UserRepositoryOption {
	return func(f *FakeUserRepository) {
		f.FGetAllUsers = append(f.FGetAllUsers, fn...)
	}
}

func OnFindUserByID(fn ...FindUserByIDFunc) UserRepositoryOption {
	return func(f *FakeUserRepository) {
		f.FFindUserByID = append(f.FFindUserByID, fn...)
	}
}

func OnCreateUser(fn ...CreateUserFunc) UserRepositoryOption {
	return func(f *FakeUserRepository) {
		f.FCreateUser = append(f.FCreateUser, fn...)
	}
}

func OnDeleteUserByID(fn ...DeleteUserByIDFunc) UserRepositoryOption {
	return func(f *FakeUserRepository) {
		f.FDeleteUserByID = append(f.FDeleteUserByID, fn...)
	}
}

func OnUpdateUser(fn ...UpdateUserFunc) UserRepositoryOption {
	return func(f *FakeUserRepository) {
		f.FUpdateUser = append(f.FUpdateUser, fn...)
	}
}

func OnExecute(fn ...ExecuteFunc) UserRepositoryOption {
	return func(f *FakeUserRepository) {
		f.FExecute = append(f.FExecute, fn...)
	}
}

func OnFindUserByName(fn ...FindUserByNameFunc) UserRepositoryOption {
	return func(f *FakeUserRepository) {
		f.FFindUserByName = append(f.FFindUserByName, fn...)
	}
}

func (f *FakeUserRepository) OnGetAllUsers(fns ...GetAllUsersFunc) {
	for _, fn := range fns {
		f.FGetAllUsers = append(f.FGetAllUsers, fn)
	}
}

func (f *FakeUserRepository) OnFindUserByID(fns ...FindUserByIDFunc) {
	for _, fn := range fns {
		f.FFindUserByID = append(f.FFindUserByID, fn)
	}
}

func (f *FakeUserRepository) OnCreateUser(fns ...CreateUserFunc) {
	for _, fn := range fns {
		f.FCreateUser = append(f.FCreateUser, fn)
	}
}

func (f *FakeUserRepository) OnDeleteUserByID(fns ...DeleteUserByIDFunc) {
	for _, fn := range fns {
		f.FDeleteUserByID = append(f.FDeleteUserByID, fn)
	}
}

func (f *FakeUserRepository) OnUpdateUser(fns ...UpdateUserFunc) {
	for _, fn := range fns {
		f.FUpdateUser = append(f.FUpdateUser, fn)
	}
}

func (f *FakeUserRepository) OnExecute(fns ...ExecuteFunc) {
	for _, fn := range fns {
		f.FExecute = append(f.FExecute, fn)
	}
}

func (f *FakeUserRepository) OnFindUserByName(fns ...FindUserByNameFunc) {
	for _, fn := range fns {
		f.FFindUserByName = append(f.FFindUserByName, fn)
	}
}

func NewFakeUserRepository(t *testing.T, opts ...UserRepositoryOption) *FakeUserRepository {
	f := &FakeUserRepository{t: t}
	for _, opt := range opts {
		opt(f)
	}
	t.Cleanup(func() {
		if f.GetAllUsersCount != len(f.FGetAllUsers) {
			t.Fatalf("expected GetAllUsers to be called %d times but got %d", len(f.FGetAllUsers), f.GetAllUsersCount)
		}
		if f.FindUserByIDCount != len(f.FFindUserByID) {
			t.Fatalf("expected FindUserByID to be called %d times but got %d", len(f.FFindUserByID), f.FindUserByIDCount)
		}
		if f.CreateUserCount != len(f.FCreateUser) {
			t.Fatalf("expected CreateUser to be called %d times but got %d", len(f.FCreateUser), f.CreateUserCount)
		}
		if f.DeleteUserByIDCount != len(f.FDeleteUserByID) {
			t.Fatalf("expected DeleteUserByID to be called %d times but got %d", len(f.FDeleteUserByID), f.DeleteUserByIDCount)
		}
		if f.UpdateUserCount != len(f.FUpdateUser) {
			t.Fatalf("expected UpdateUser to be called %d times but got %d", len(f.FUpdateUser), f.UpdateUserCount)
		}
		if f.ExecuteCount != len(f.FExecute) {
			t.Fatalf("expected Execute to be called %d times but got %d", len(f.FExecute), f.ExecuteCount)
		}
		if f.FindUserByNameCount != len(f.FFindUserByName) {
			t.Fatalf("expected FindUserByName to be called %d times but got %d", len(f.FFindUserByName), f.FindUserByNameCount)
		}
	})
	return f
}

func (f *FakeUserRepository) GetAllUsers() ([]User, error) {
	var idx = f.GetAllUsersCount
	if f.GetAllUsersCount >= len(f.FGetAllUsers) {
		idx = len(f.FGetAllUsers) - 1
	}
	if len(f.FGetAllUsers) != 0 {
		o1, o2 := f.FGetAllUsers[idx]()
		f.GetAllUsersCount++
		return o1, o2
	}
	return nil, nil
}

func (f *FakeUserRepository) FindUserByID(id int) (u User, err error) {
	var idx = f.FindUserByIDCount
	if f.FindUserByIDCount >= len(f.FFindUserByID) {
		idx = len(f.FFindUserByID) - 1
	}
	if len(f.FFindUserByID) != 0 {
		u, err := f.FFindUserByID[idx](id)
		f.FindUserByIDCount++
		return u, err
	}
	return User{}, nil
}

func (f *FakeUserRepository) CreateUser(user User) error {
	var idx = f.CreateUserCount
	if f.CreateUserCount >= len(f.FCreateUser) {
		idx = len(f.FCreateUser) - 1
	}
	if len(f.FCreateUser) != 0 {
		o1 := f.FCreateUser[idx](user)
		f.CreateUserCount++
		return o1
	}
	return nil
}

func (f *FakeUserRepository) DeleteUserByID(id int) error {
	var idx = f.DeleteUserByIDCount
	if f.DeleteUserByIDCount >= len(f.FDeleteUserByID) {
		idx = len(f.FDeleteUserByID) - 1
	}
	if len(f.FDeleteUserByID) != 0 {
		o1 := f.FDeleteUserByID[idx](id)
		f.DeleteUserByIDCount++
		return o1
	}
	return nil
}

func (f *FakeUserRepository) UpdateUser(old, new User) error {
	var idx = f.UpdateUserCount
	if f.UpdateUserCount >= len(f.FUpdateUser) {
		idx = len(f.FUpdateUser) - 1
	}
	if len(f.FUpdateUser) != 0 {
		o1 := f.FUpdateUser[idx](old, new)
		f.UpdateUserCount++
		return o1
	}
	return nil
}

func (f *FakeUserRepository) Execute(user User) error {
	var idx = f.ExecuteCount
	if f.ExecuteCount >= len(f.FExecute) {
		idx = len(f.FExecute) - 1
	}
	if len(f.FExecute) != 0 {
		o1 := f.FExecute[idx](user)
		f.ExecuteCount++
		return o1
	}
	return nil
}

func (f *FakeUserRepository) FindUserByName(name string) (User, error) {
	var idx = f.FindUserByNameCount
	if f.FindUserByNameCount >= len(f.FFindUserByName) {
		idx = len(f.FFindUserByName) - 1
	}
	if len(f.FFindUserByName) != 0 {
		o1, o2 := f.FFindUserByName[idx](name)
		f.FindUserByNameCount++
		return o1, o2
	}
	return User{}, nil
}
