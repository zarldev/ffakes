For Fakes Sake
fake it ~~til~~ after you make it
Everyone loves the implicit nature of interfaces in Go it just becomes second nature and makes life as a Go developer a little different when writing composable software.  It gives us a level of freedom, for example, pulling in a 3rd party library realizing you are using 5% of the features, and you want to do some testing - so you just write an interface for the methods you use and *go* on with your day.

The next thing we want to do is write some tests so we use our new interface and generate a mock using that industry standard frame work - [mockery](https://github.com/vektra/mockery).

I like to use inline go:generate flags so I can automate some steps along with the build process so the command for generating the mock is:

```
//go:generate mockery --name=UserRepository --filename=repository_mock.go --inpackage
```

After I run `go generate ./...` the `repository_mock.go` file is generated and we can now substitute the dependency of our SUT (*Subject Under Test*). With our new `MockUserRepository` we can now mock some function calls for the dependency.

Lets start with the structure of the underlying interface the `UserRepository` interface is a standard Repository interface for storing the `User` and looks like the following:

```golang
type User struct {
	ID   int
	Name string
}

type UserRepository interface {
	FindUserByID(id int) (u User, err error)
	CreateUser(user User) error
	DeleteUserByID(id int) error
	UpdateUser(old, new User) error
	Execute(user User) error
	FindUserByName(name string) (User, error)
}
```
So using this mock we have to start with the setup of our expectations for the `CreateUser` method; This is simple as it takes a struct of our creation and just returns an error. As per good testing, let's produce a success scenario and a failing scenario.

### Successful
```golang
// Instantiate new mock
mockRepo := repository.NewMockUserRepository(t)

// Set Expectations
// Explicit
mockRepo.On("CreateUser", repository.User{ID: 1, Name: "test1"}).Return(nil)
// Execution
u := repository.User{
    ID: 1, 
    Name: "test1"
}

err := mockRepo.CreateUser(u)
if err != nil {
    t.Fatal("test failed")
}
```

### Unsuccessful
```golang
// Instantiate new mock
mockRepo := repository.NewMockUserRepository(t)

// Set Expectations
// Generalised 
mockRepo.On("CreateUser", mock.AnythingOfType("repository.User")).Return(fmt.Errorf("expected"))

// Execution
u := repository.User{
    ID: 1, 
    Name: "test1"
}

err := mockRepo.CreateUser(u)
if err != nil {
    t.Fatal("test failed")
}
```

Note the difference in the passing of the arguments to the `mockRepo.On` function call.  Lets break this down:

1. The method signature for the expectations does not match the interface we have mocked rather a description of what we mocked.
2. As such the `mockRepo.On` method takes a `string` of the method name we want to mock as the 1st parameter.
3. The 2nd parameter of the `mockRepo.On` method is the variadic argument list for the method named in the 1st parameter.  
4. This argument list is 1 or more of type `any`.
5. The argument `mock.AnythingOfType("repository.User")` is again using a string but this time as a representation of a type.

Using this to me seems to have missed the point of the implicit nature of Go's interfaces.  I feel that using strings as the names of types or methods is clumsy and prone to errors.  They are also refactor proof, if I refactor for example the `CreateUser` to `StoreUser` and regenerate my mocks, the tests are still compilable but reference a function that does not exist.  The variadic nature to be the supplied parameters is also prone to errors; again its refactor proof but also allows me to construct compilable code that does not match the parameter list of the method named in the 1st parameter.  I can do this in a multitude of ways

* by providing more or less arguments than the method expects.
* pointer semantics - having to make sure the string has the `*` or `&` prepended to the string.
* when using the `mock.AnythingOfType` making sure the imports used on the type are correct.

```golang
// Actual
mockRepo.On("CreateUser", mock.AnythingOfType("repository.User")).Return(nil)
// Whatever you want it to be
mockRepo.On("IDoNotExist", struct{}{}, os.Stderr, nil, -1).Return(0,true,false)
```
The returned test error in this case is `The code you are testing needs to make 2 more call(s).` as the mock has expectations that both`CreateUser`and `IDoNotExist` to be called however the mock has no way to know that `IDoNotExist` in fact does not exist and hence will NEVER be called.

All in all this is a really flexible approach and allows some of the features that some people really rely on in there mocks.  However in practice in the real world I find it really ends up a bit of a hell-scape of shitty copy pasted code of the last previously working test.  And as the DSL ends up taking its tole on developers will to live with the amount of hoops that are required to be jumped through.   But finally you have jumped though all the hoops and have working tests and everything is rosy. UNTIL you then decide to do some refactor work and all your tests blow up because everything still compiles even if the interface has changed. Now its a hunt for that magic string that has an old method name or wrong pointer semantics.  In my experience this leads to a proliferation of the weakest abstraction available `mock.Anything`, not even `mock.AnythingOfType` just anything... Now these tests hold close to 0 value and actually do not test much at all apart from your patience.

This is not to say that I don't acknowledge the use-case for mocks or when they really do provide an easier way to improve the ease of writing correct tests. There is also the added benefit that you don't have to write any of the implementation details or understand the details of what the interface does.  They can also make it easier to facilitate complex scenarios of the same method call.

What I find myself doing a lot of the time I'm in this situation is just go ahead and generate a `fake` with some super simple implementation of the required interface.  If you have followed good practice your interfaces should not have lots of methods in them ~~like my `UserRepository` example~~.  But for the fake it stores a slice of functions for each method signature in the defined interface.  I can then override those as required to return the expected result.  So the `FakeUserRepository` is define as:

```golang
// Generated fake struct that conforms to the UserRepository interface 
type FakeUserRepository struct {
	// Each entry is a slice of functions that conform to each method
	FFindUserByID []func(id int) (u User, err error)
	FCreateUser []func(user User) error
	FDeleteUserByID []func(id int) error
	FUpdateUser []func(old, new User) error
	FExecute []func(user User) error
	FFindUserByName []func(name string) (User, error)
}
```
I can then just create the fake and prescribe the behaviour on top of the existing function calls. I usually add some `On` methods to allow easy additional expectations for each method.

```golang
func (f *UserRepositoryFake) OnCreateUser(fns ...CreateUserFunc) {
	for _, fn := range fns {
		f.FCreateUser = append(f.FCreateUser, fn)
	}
}
```

These can then be used to set expectations and implement the call simply.

```golang
repo.OnCreateUser(func(user repository.User) error {
		return nil
	})
```

Or for example setup test data then use it within the fake but to implement behaviour.

```golang
userIDMap := make(map[int]repository.User, 0)
userNameMap := make(map[string]repository.User, 0)

repo.OnCreateUser(func(user repository.User) error {
		userIDMap[user.ID] = user
		userNameMap[user.Name] = user
		return nil
	})
```
I find the fact that it conforms to the interface definitions allows to focus on the inputs and outputs for the test scenario.  Some people will say just go write a `InMemoryUserRepository` and use that as the fake for all your calls.  This is a valid comment and in some case that should be the correct call to make.  It does mean that you have to also write some of the common features people like to rely on in tests, validating the method has been called the correct amount of times.  Or setting expectations for methods to be called and then them not to be or the inverse when the method is called too many times.  We can achieve this by adding a `CreateUserCount` of type `int` and incrementing it when the method is called.  The other is a little trickier to do, we can however leverage the `testing.T` types `Cleanup` method which is ran after each test, here we can build or expectation checks around method call counts, where we check the count against the number of functions in our corresponding method slice.

```golang
t.Cleanup(func() {
		if f.CreateUserCount != len(f.FCreateUser) {
			t.Fatalf("expected CreateUser to be called %d times but got %d", len(f.FCreateUser), f.CreateUserCount)
		}
})
```

And **YES** I am fully aware that at this point the line between fake and mock is well and truly gone and some would argue this is more mock than fake, but at this point I just want an easy abstraction for writing tests you can call it what you want 😂.

I have found myself doing things like this for building out test case so many times, and with the understanding of the `ast` libraries in Go I decided to write a generator for this. It's called `ffakes` a play on the classic ffs acronym and it will build out fakes with all the helpers for testing.

With a little `go:generate` magic `//go:generate ffakes -i UserRepository` the `FakeUserRepository` is generated.  It has the helper stuff explained out before but also has some other helpers, the `NewFakeUserRepository` function take the `testing.T` and sets up the `Cleanup` call it also generates `Option` functions for configuring the `FakeUserRepository` on instantiation and also concrete function definitions.

```golang
type UserRepositoryOption func(f *FakeUserRepository)

func OnCreateUser(fn ...CreateUserFunc) UserRepositoryOption {
	return func(f *FakeUserRepository) {
		f.FCreateUser = append(f.FCreateUser, fn...)
	}
}

repo := repository.NewFakeUserRepository(t,
	repository.OnCreateUser(func(user repository.User) error {
		return nil
	}),
)
```

The defined functions allow us to easily create the function as per our fake but then set expectations on number of calls for each method.

```golang
f := repository.FindUserByIDFunc(func(id int) (u repository.User, err error) {
	user, ok := userIDMap[id]
	if !ok {
		return repository.User{}, fmt.Errorf("unexpected user id")
	}
	return user, nil
})
// Setup
repo.OnFindUserByID(f, f, f)
```

So now we have a type safe way to quickly generate a test fake/mock and configure it in a type safe way with all the usual expectation checks that are relied upon in many code bases.

The `ffakes` project is available on Github at [ffakes](https://www.github.com/zarldev/ffakes)
