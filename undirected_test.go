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

type UndirectedMutableGraphSuite struct {
	MutableGraphSuite
}

type UndirectedGraphSuite struct {
	GraphSuite
}

var _ = Suite(&UndirectedMutableGraphSuite{
	MutableGraphSuite{Factory: u_fact},
})

var _ = Suite(&UndirectedGraphSuite{
	GraphSuite{Factory: u_fact},
})
