# ffakes

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![build](https://github.com/zarldev/ffakes/actions/workflows/go.yml/badge.svg)

# Introduction 
ffakes is a tool to generate fakes for interfaces in Go code. It generates a new file with the same name as the original file with "_fakes" appended to the name. The generated file contains a struct with fields for each method in the interface where each field is a slice of functions that can be used to fake the series of calls to the function. The generated file also contains a type alias for each function type and a set of options to configure the fake.  The fake struct can also be configured with the same options after initialization. The generated file is placed in the same directory as the original file by default but can be overridden with the output flag. If there is more calls to a method than the number of functions in the slice, while being used for testing the testing.T.Fatal function is called for the *testing.T object required for the New method.

# Installation 
```
go install github.com/zarldev/ffakes@latest
```

# Usage 
```
$ ffakes -h

