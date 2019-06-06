package gogl

import (
	"gopkg.in/fatih/set.v0"
)

// Shared helper function for edge lists to enumerate vertices.
func elVertices(el interface{}, fn VertexStep) {
	set := set.New(set.NonThreadSafe)

	el.(EdgeEnumerator).Edges(func(e Edge) (terminate bool) {
		set.Add(e.Both())
		return
	})

	for _, v := range set.List() {
		if fn(v) {
			return
		}
	}
}

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

func (el EdgeList) Vertices(fn VertexStep) {
	elVertices(el, fn)
}

func (el EdgeList) Edges(fn EdgeStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

// An ArcList is a naive DigraphSource implementation that is backed only by an arc slice.
type ArcList []Arc

func (el ArcList) Vertices(fn VertexStep) {
	elVertices(el, fn)
}

func (el ArcList) Edges(fn EdgeStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

func (el ArcList) Arcs(fn ArcStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

// A WeightedEdgeList is a naive GraphSource implementation that is backed only by an edge slice.
//
// This variant is for weighted edges.
type WeightedEdgeList []WeightedEdge

func (el WeightedEdgeList) Vertices(fn VertexStep) {
	elVertices(el, fn)
}

func (el WeightedEdgeList) Edges(fn EdgeStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

// A WeightedArcList is a naive DigraphSource implementation that is backed only by an arc slice.
type WeightedArcList []Arc

func (el WeightedArcList) Vertices(fn VertexStep) {
	elVertices(el, fn)
}

func (el WeightedArcList) Edges(fn EdgeStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

func (el WeightedArcList) Arcs(fn ArcStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

// A LabeledEdgeList is a naive GraphSource implementation that is backed only by an edge slice.
//
// This variant is for labeled edges.
type LabeledEdgeList []LabeledEdge

func (el LabeledEdgeList) Vertices(fn VertexStep) {
	elVertices(el, fn)
}

func (el LabeledEdgeList) Edges(fn EdgeStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

// A LabeledArcList is a naive DigraphSource implementation that is backed only by an arc slice.
type LabeledArcList []Arc

func (el LabeledArcList) Vertices(fn VertexStep) {
	elVertices(el, fn)
}

func (el LabeledArcList) Edges(fn EdgeStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

func (el LabeledArcList) Arcs(fn ArcStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

// A DataEdgeList is a naive GraphSource implementation that is backed only by an edge slice.
//
// This variant is for labeled edges.
type DataEdgeList []DataEdge

func (el DataEdgeList) Vertices(fn VertexStep) {
	elVertices(el, fn)
}

func (el DataEdgeList) Edges(fn EdgeStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

// A DataArcList is a naive DigraphSource implementation that is backed only by an arc slice.
type DataArcList []Arc

func (el DataArcList) Vertices(fn VertexStep) {
	elVertices(el, fn)
}

func (el DataArcList) Edges(fn EdgeStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}

func (el DataArcList) Arcs(fn ArcStep) {
	for _, e := range el {
		if fn(e) {
			return
		}
	}
}
