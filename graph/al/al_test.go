package al

import (
	"testing"

	"github.com/sdboyer/gogl/spec"
	"github.com/sdboyer/gocheck"
)

// Hook gocheck into the go test runner
func TestHookup(t *testing.T) { gocheck.TestingT(t) }

func init() {
	for gp, _ := range alCreators {
		spec.SetUpTestsFromSpec(gp, AdjacencyList)
	}
}

