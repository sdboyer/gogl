package gogl

import (
	"fmt"
	. "launchpad.net/gocheck"
)

var _ = fmt.Println

var u_fact = &GraphFactory{
	CreateMutableGraph: func() MutableGraph {
		return NewUndirected()
	},
	CreateGraph: func(edges []Edge) Graph {
		return NewUndirectedFromEdgeSet(edges)
	},
}

var _ = Suite(&MutableGraphSuite{
	Factory: u_fact,
})

var _ = Suite(&GraphSuite{
	Factory: u_fact,
})
