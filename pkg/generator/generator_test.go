package generator_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/zarldev/ffakes/pkg/generator"
)

var (
	tests = []struct {
		inputFile   string
		outputFile  string
		interfaces  []string
		packageName string

		shouldError   bool
		expectedError string
	}{
		{
			inputFile:   "tests/repository/repository.go",
			outputFile:  "tests/repository/repository_fakes.go",
			interfaces:  []string{"UserRepository"},
			packageName: "repository",
		},
		{
			inputFile:   "tests/pubsub/pubsub.go",
			outputFile:  "tests/pubsub/pubsub_fakes.go",
			interfaces:  []string{"Publisher", "Subscriber", "Broker"},
			packageName: "repository",
		},
		{
			inputFile:     "tests/invalid.go",
			outputFile:    "tests/invalid_fakes.go",
			interfaces:    []string{"InvalidInterface"},
			packageName:   "tests/invalid",
			shouldError:   true,
			expectedError: "failed to parse file while generating fakes: open tests/invalid.go: no such file or directory",
		},
	}
)

func TestMain(m *testing.M) {
	m.Run()
}

func cleanUp() {
	for _, test := range tests {
		if test.outputFile != "" {
			err := os.Remove(test.outputFile)
			if err != nil {
				slog.Error("failed to remove file", slog.String("file", test.outputFile), slog.String("error", err.Error()))
			}
		}
	}
}

func TestParseAndGenerate(t *testing.T) {
	cleanUp()
	t.Run("when parsing and generating for a file", func(t *testing.T) {
		for _, test := range tests {
			t.Run(fmt.Sprintf("should return the package name and interfaces for %s", test.inputFile), func(t *testing.T) {
				err := generator.ParseAndGenerate(test.inputFile, test.interfaces, test.packageName)
				if test.shouldError {
					if err == nil {
						t.Errorf("expected an error but got nil")
					}
					if err.Error() != test.expectedError {
						t.Errorf("expected error to be '%s' but got '%s'", test.expectedError, err.Error())
					}
					return
				}
				if err != nil {
					t.Errorf("expected no error but got %v", err)
				}
				if _, err := os.Stat(test.outputFile); os.IsNotExist(err) {
					t.Errorf("expected file '%s' to be created but it was not", test.outputFile)
				}
			})
		}
	})
}
