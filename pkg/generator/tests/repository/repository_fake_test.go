package repository_test

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/zarldev/ffakes/pkg/generator/tests/repository"
// )

// func TestRepositoryFake(t *testing.T) {
// 	user1 := repository.User{
// 		ID:   1,
// 		Name: "test1",
// 	}
// 	user2 := repository.User{
// 		ID:   2,
// 		Name: "test2",
// 	}
// 	user3 := repository.User{
// 		ID:   3,
// 		Name: "test3",
// 	}
// 	userIDMap := map[int]repository.User{
// 		user1.ID: user1,
// 		user2.ID: user2,
// 		user3.ID: user3,
// 	}
// 	userNameMap := map[string]repository.User{
// 		user1.Name: user1,
// 		user2.Name: user2,
// 		user3.Name: user3,
// 	}
// 	repo := repository.NewFakeUserRepository(t,
// 		repository.OnCreateUser(func(user repository.User) error {
// 			return nil
// 		}),
// 	)
// 	t.Run("when finding a user by id", func(t *testing.T) {
// 		f := repository.FindUserByIDFunc(func(id int) (u repository.User, err error) {
// 			user, ok := userIDMap[id]
// 			if !ok {
// 				return repository.User{}, fmt.Errorf("unexpected user id")
// 			}
// 			return user, nil
// 		})
// 		// Setup
// 		repo.OnFindUserByID(f, f, f)

// 		// Test
// 		u, err := repo.FindUserByID(1)
// 		if err != nil {
// 			t.Errorf("expected no error but got %v", err)
// 		}
// 		if u != user1 {
// 			t.Errorf("expected user to be %v but got %v", user1, u)
// 		}

// 		u, err = repo.FindUserByID(2)
// 		if err != nil {
// 			t.Errorf("expected no error but got %v", err)
// 		}
// 		if u != user2 {
// 			t.Errorf("expected user to be %v but got %v", user2, u)
// 		}

// 		u, err = repo.FindUserByID(3)
// 		if err != nil {
// 			t.Errorf("expected no error but got %v", err)
// 		}
// 		if u != user3 {
// 			t.Errorf("expected user to be %v but got %v", user3, u)
// 		}

// 	})
// 	t.Run("when creating a user", func(t *testing.T) {
// 		// Setup
// 		repo.OnCreateUser(func(user repository.User) error {
// 			return nil
// 		})
// 		// Test
// 		err := repo.CreateUser(user1)
// 		if err != nil {
// 			t.Errorf("expected no error but got %v", err)
// 		}
// 	})
// 	t.Run("when deleting a user by id", func(t *testing.T) {
// 		// Setup
// 		repo.OnDeleteUserByID(func(id int) error {
// 			return nil
// 		})
// 		// Test
// 		err := repo.DeleteUserByID(1)
// 		if err != nil {
// 			t.Errorf("expected no error but got %v", err)
// 		}
// 	})
// 	t.Run("when updating a user", func(t *testing.T) {
// 		// Setup
// 		repo.OnUpdateUser(func(old repository.User, new repository.User) error {
// 			return nil
// 		})
// 		// Test
// 		err := repo.UpdateUser(user1, user2)
// 		if err != nil {
// 			t.Errorf("expected no error but got %v", err)
// 		}
// 	})
// 	t.Run("when executing a user", func(t *testing.T) {
// 		// Setup
// 		repo.OnExecute(func(user repository.User) error {
// 			return nil
// 		})
// 		// Test
// 		err := repo.Execute(user1)
// 		if err != nil {
// 			t.Errorf("expected no error but got %v", err)
// 		}
// 	})
// 	t.Run("when finding a user by name", func(t *testing.T) {
// 		// Setup
// 		repo.OnFindUserByName(func(name string) (u repository.User, err error) {
// 			user, ok := userNameMap[name]
// 			if !ok {
// 				return repository.User{}, fmt.Errorf("unexpected user name")
// 			}
// 			return user, nil
// 		})
// 		// Test
// 		u, err := repo.FindUserByName("test1")
// 		if err != nil {
// 			t.Errorf("expected no error but got %v", err)
// 		}
// 		if u != user1 {
// 			t.Errorf("expected user to be %v but got %v", user1, u)
// 		}
// 	})
// }
