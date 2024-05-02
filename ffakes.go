// ffakes is a tool to generate fakes for interfaces in Go code.
// It generates a new file with the same name as the original file with "_fakes" appended to the name.
// The generated file contains a struct with fields for each method in the interface
// where each field is a slice of functions that can be used to fake the series of calls to the function.
// The generated file also contains a type alias for each function type and a set of options to configure the fake.
// The fake struct can also be configured with the same options after initialization. The generated
// file is placed in the same directory as the original file by default but can be overridden with the output flag.
// If there is more calls to a method than the number of functions in the slice, while being used for
// testing the testing.T.Fatal function is called for the *testing.T struct required for the New method
//
// Usage: ffakes [options] filename
//
// Options:
//
// ffakes [options] -i, --interfaces string Comma separated list of interfaces to fake
//
// This can also be used in a go generate directive.
// Example:
// //go:generate ffakes -i UserRepository
//
// This will generate a new file called repository_fakes.go in the same directory as the original file.
// The generated file will contain a struct called UserRepositoryFake with fields for each method in the UserRepository interface.
//
// https://www.zarl.dev
package main

import (
	"github.com/zarldev/ffakes/pkg/app"
)

func main() {
	app.Run()
}
