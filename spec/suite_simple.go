package spec

import (
	"fmt"
	"math"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

/* SimpleGraphSuite - tests for simple graph methods */

type SimpleGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *SimpleGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *SimpleGraphSuite) TestDensity(c *C) {
	c.Assert(math.IsNaN(s.Factory(NullGraph).(SimpleGraph).Density()), Equals, true)

	g := s.Factory(GraphFixtures["pair"]).(SimpleGraph)
	if s.Directed {
		c.Assert(g.Density(), Equals, float64(0.5))
	} else {
		c.Assert(g.Density(), Equals, float64(1))
	}

	g = s.Factory(GraphFixtures["2e3v"]).(SimpleGraph)
	if s.Directed {
		c.Assert(g.Density(), Equals, float64(2)/float64(6))
	} else {
		c.Assert(g.Density(), Equals, float64(2)/float64(3))
	}
}

