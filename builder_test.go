package gogl

import (
	. "github.com/sdboyer/gocheck"
)

type SpecTestSuite struct{}

var _ = Suite(&SpecTestSuite{})

const baseline = G_UNDIRECTED | G_SIMPLE | G_BASIC | G_MUTABLE

func (s *SpecTestSuite) TestInitialState(c *C) {
	c.Assert(Spec().Props == baseline, Equals, true)
}

// Because as long as we're playing the coverage game, why not be ridiculous?
func (s *SpecTestSuite) permuteField() []GraphProperties {
	field := make([]GraphProperties, 16, 16)
	for i := uint(0); i < 16; i++ {
		field[i] = 1 << i
	}

	return field
}

func (s *SpecTestSuite) TestMutators(c *C) {
	var spec GraphSpec
	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.Directed().Props&G_DIRECTED == G_DIRECTED, Equals, true)
		c.Assert(spec.Directed().Props&G_UNDIRECTED == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.Undirected().Props&G_UNDIRECTED == G_UNDIRECTED, Equals, true)
		c.Assert(spec.Undirected().Props&G_DIRECTED == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.Basic().Props&G_BASIC == G_BASIC, Equals, true)
		c.Assert(spec.Basic().Props&(G_LABELED|G_WEIGHTED|G_DATA) == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.Labeled().Props&G_LABELED == G_LABELED, Equals, true)
		c.Assert(spec.Labeled().Props&G_BASIC == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.Weighted().Props&G_WEIGHTED == G_WEIGHTED, Equals, true)
		c.Assert(spec.Weighted().Props&G_BASIC == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.DataEdges().Props&G_DATA == G_DATA, Equals, true)
		c.Assert(spec.DataEdges().Props&G_BASIC == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.SimpleGraph().Props&G_SIMPLE == G_SIMPLE, Equals, true)
		c.Assert(spec.SimpleGraph().Props&(G_LOOPS|G_PARALLEL) == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.MultiGraph().Props&G_PARALLEL == G_PARALLEL, Equals, true)
		c.Assert(spec.MultiGraph().Props&(G_LOOPS|G_SIMPLE) == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.PseudoGraph().Props&(G_LOOPS|G_PARALLEL) == (G_LOOPS|G_PARALLEL), Equals, true)
		c.Assert(spec.PseudoGraph().Props&G_SIMPLE == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.Parallel().Props&G_PARALLEL == G_PARALLEL, Equals, true)
		c.Assert(spec.Parallel().Props&G_SIMPLE == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.Loop().Props&G_LOOPS == G_LOOPS, Equals, true)
		c.Assert(spec.Loop().Props&G_SIMPLE == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.Mutable().Props&G_MUTABLE == G_MUTABLE, Equals, true)
		c.Assert(spec.Mutable().Props&G_IMMUTABLE == 0, Equals, true)
	}

	for _, spec.Props = range s.permuteField() {
		c.Assert(spec.Immutable().Props&G_IMMUTABLE == G_IMMUTABLE, Equals, true)
		c.Assert(spec.Immutable().Props&G_MUTABLE == 0, Equals, true)
	}
}
