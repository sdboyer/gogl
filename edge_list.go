package gogl

import (
	"gopkg.in/fatih/set.v0"
)

// An EdgeList is a naive GraphSource implementation that is backed only by an edge slice.
//
// EdgeLists are primarily intended for use as fixtures.
//
// It is inherently impossible for an EdgeList to represent a vertex isolate (degree 0) fully
// correctly, as vertices are only described in the context of their edges. One can be a
// little hacky, though, and represent one with a loop. As gogl expects graph implementations
// to simply discard loops if they are disallowed by the graph's constraints (i.e., in simple
// and multigraphs), they *should* be interpreted as vertex isolates.
type EdgeList []Edge

func (el EdgeList) EachVertex(fn VertexLambda) {
	for _, v := range CollectVertices(el) {
		if fn(v) {
			return
		}
	}
}

func (el EdgeList) EachEdge(fn EdgeLambda) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

func (el EdgeList) Order() int {
	set := set.NewNonTS()

	el.EachEdge(func(e Edge) (terminate bool) {
		set.Add(e.Both())
		return
	})

	return set.Size()
}

