package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

/* Counting suites - tests for Size() and Order() */

type OrderSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *OrderSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *OrderSuite) TestOrder(c *C) {
	c.Assert(s.Factory(NullGraph).(VertexCounter).Order(), Equals, 0)
	c.Assert(s.Factory(GraphFixtures["2e3v"]).(VertexCounter).Order(), Equals, 3)
}

type SizeSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *SizeSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *SizeSuite) TestSize(c *C) {
	c.Assert(s.Factory(NullGraph).(EdgeCounter).Size(), Equals, 0)
	c.Assert(s.Factory(GraphFixtures["2e3v"]).(EdgeCounter).Size(), Equals, 2)
}
