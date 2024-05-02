package repository_test

import (
	"testing"

	"github.com/zarldev/ffakes/pkg/generator"
)

func TestParseAndGenerate(t *testing.T) {
	t.Run("when parsing and generating for a file", func(t *testing.T) {
		t.Run("should return the package name and interfaces", func(t *testing.T) {
			err := generator.ParseAndGenerate("repository.go", []string{"UserRepository", "Publisher", "Subscriber", "Broker"}, "")
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
		})

	})
}
