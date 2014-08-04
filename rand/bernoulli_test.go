package rand

import (
	"fmt"
	stdrand "math/rand"
	"testing"
	"time"

	. "github.com/sdboyer/gocheck"
	"github.com/sdboyer/gogl"
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
		"dir_stable":   BernoulliDistribution(10, 0.5, true, true, r),
		"und_stable":   BernoulliDistribution(10, 0.5, false, true, r),
		"dir_unstable": BernoulliDistribution(10, 0.5, true, false, r),
		"und_unstable": BernoulliDistribution(10, 0.5, false, false, r),
	}
}

func (s *BernoulliTest) TestLengthChecks(c *C) {
	c.Assert(gogl.Order(s.graphs["dir_stable"]), Equals, 10)
	c.Assert(gogl.Order(s.graphs["und_stable"]), Equals, 10)
	c.Assert(gogl.Order(s.graphs["dir_unstable"]), Equals, 10)
	c.Assert(gogl.Order(s.graphs["und_unstable"]), Equals, 10)
}

func (s BernoulliTest) TestEachVertex(c *C) {
	sl := make([]int, 0, 40)

	for _, g := range s.graphs {
		g.EachVertex(func(v gogl.Vertex) (terminate bool) {
			sl = append(sl, v.(int))
			return
		})
	}

	c.Assert(len(sl), Equals, 40)

	for k, v := range sl {
		c.Assert(k%10, Equals, v)
	}
}
