package rand

import (
	"fmt"
	stdrand "math/rand"
	"testing"
	"time"

	. "github.com/sdboyer/gocheck"
	"github.com/sdboyer/gogl"
	"gopkg.in/fatih/set.v0"
)

var fml = fmt.Println

func TestRand(t *testing.T) { TestingT(t) }

type BernoulliTest struct {
	graphs map[string]gogl.GraphSource
}

var _ = Suite(&BernoulliTest{})

func (s *BernoulliTest) SetUpSuite(c *C) {
	r := stdrand.NewSource(time.Now().UnixNano())
	s.graphs = map[string]gogl.GraphSource{
		"dir_stable":         BernoulliDistribution(10, 0.5, true, true, r),
		"und_stable":         BernoulliDistribution(10, 0.5, false, true, r),
		"dir_unstable":       BernoulliDistribution(10, 0.5, true, false, r),
		"und_unstable":       BernoulliDistribution(10, 0.5, false, false, r),
		"und_unstable_nosrc": BernoulliDistribution(10, 0.5, false, false, nil),
	}
}

func (s *BernoulliTest) TestLengthChecks(c *C) {
	c.Assert(gogl.Order(s.graphs["dir_stable"]), Equals, 10)
	c.Assert(gogl.Order(s.graphs["und_stable"]), Equals, 10)
	c.Assert(gogl.Order(s.graphs["dir_unstable"]), Equals, 10)
	c.Assert(gogl.Order(s.graphs["und_unstable"]), Equals, 10)
	c.Assert(gogl.Order(s.graphs["und_unstable_nosrc"]), Equals, 10)
}

func (s *BernoulliTest) TestProbabilityRange(c *C) {
	f1 := func() {
		BernoulliDistribution(1, -0.0000001, true, true, nil)
	}

	f2 := func() {
		BernoulliDistribution(1, 1.0, true, true, nil)
	}
	c.Assert(f1, PanicMatches, "ρ must be in the range \\[0\\.0,1\\.0\\).")
	c.Assert(f2, PanicMatches, "ρ must be in the range \\[0\\.0,1\\.0\\).")
}

func (s *BernoulliTest) TestVertices(c *C) {
	sl := make([]int, 0, 50)

	for _, g := range s.graphs {
		g.Vertices(func(v gogl.Vertex) (terminate bool) {
			sl = append(sl, v.(int))
			return
		})

	}

	c.Assert(len(sl), Equals, 50)

	for k, v := range sl {
		c.Assert(k%10, Equals, v)
	}
}

func (s *BernoulliTest) TestVerticesTermination(c *C) {
	var hit int
	s.graphs["dir_stable"].Vertices(func(v gogl.Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)

	s.graphs["dir_unstable"].Vertices(func(v gogl.Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 2)
}

func (s *BernoulliTest) TestEachEdgeCount(c *C) {
	// Given that this is a rand count, our testing options are curtailed
	for gn, g := range s.graphs {
		hit := 0
		g.EachEdge(func(e gogl.Edge) (terminate bool) {
			hit++
			return
		})

		switch gn {
		case "dir_stable", "dir_unstable":
			c.Assert(hit <= 90, Equals, true)
			c.Assert(hit >= 0, Equals, true)
		case "und_stable", "und_unstable", "und_unstable_nosrc":
			c.Assert(hit <= 45, Equals, true)
			c.Assert(hit >= 0, Equals, true)
		}
	}
}

func (s *BernoulliTest) TestEachEdgeStability(c *C) {
	setd := set.NewNonTS()
	setu := set.NewNonTS()
	var hitu, hitd int

	dg := BernoulliDistribution(10, 0.5, true, true, nil)
	dg.EachEdge(func(e gogl.Edge) (terminate bool) {
		setd.Add(e)
		return
	})

	dg.EachEdge(func(e gogl.Edge) (terminate bool) {
		c.Assert(setd.Has(e), Equals, true)
		hitd++
		return
	})

	c.Assert(setd.Size(), Equals, hitd)
	c.Assert(dg.(gogl.EdgeCounter).Size(), Equals, hitd)

	ug := BernoulliDistribution(10, 0.5, false, true, nil)
	ug.EachEdge(func(e gogl.Edge) (terminate bool) {
		setu.Add(e)
		return
	})

	ug.EachEdge(func(e gogl.Edge) (terminate bool) {
		c.Assert(setu.Has(e), Equals, true)
		hitu++
		return
	})

	c.Assert(setu.Size(), Equals, hitu)
	c.Assert(ug.(gogl.EdgeCounter).Size(), Equals, hitu)

}

func (s *BernoulliTest) TestEachEdgeTermination(c *C) {
	var hit int
	s.graphs["dir_unstable"].EachEdge(func(e gogl.Edge) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)

	s.graphs["und_unstable"].EachEdge(func(e gogl.Edge) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 2)

	gogl.CollectEdges(s.graphs["und_stable"]) // To populate the cache
	s.graphs["und_stable"].EachEdge(func(e gogl.Edge) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 3)

	gogl.CollectEdges(s.graphs["dir_stable"])
	s.graphs["dir_stable"].EachEdge(func(e gogl.Edge) bool {
		hit++
		return true
	})
	c.Assert(hit, Equals, 4)
}

func (s *BernoulliTest) TestEachArcStability(c *C) {
	setd := set.NewNonTS()
	var hitd int

	g := BernoulliDistribution(10, 0.5, true, true, nil).(gogl.DigraphSource)
	g.EachArc(func(e gogl.Arc) (terminate bool) {
		setd.Add(e)
		return
	})

	g.EachArc(func(e gogl.Arc) (terminate bool) {
		c.Assert(setd.Has(e), Equals, true)
		hitd++
		return
	})

	c.Assert(setd.Size(), Equals, hitd)
	c.Assert(g.(gogl.EdgeCounter).Size(), Equals, hitd)
}

func (s *BernoulliTest) TestEachArcTermination(c *C) {
	var hit int
	s.graphs["dir_unstable"].(gogl.DigraphSource).EachArc(func(e gogl.Arc) bool {
		hit++
		return true
	})
	c.Assert(hit, Equals, 1)

	gogl.CollectEdges(s.graphs["dir_stable"])
	s.graphs["dir_stable"].(gogl.DigraphSource).EachArc(func(e gogl.Arc) bool {
		hit++
		return true
	})
	c.Assert(hit, Equals, 2)

	s.graphs["dir_stable"].EachEdge(func(e gogl.Edge) bool {
		hit++
		return true
	})
	c.Assert(hit, Equals, 3)
}
