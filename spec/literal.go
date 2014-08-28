package spec

import (
	. "github.com/sdboyer/gogl"
)

// Define a graph literal fixture for testing here.
// The literal has two edges and four vertices; one vertex is an isolate.
//
// bool state indicates whether using a transpose or not.
type GraphLiteralFixture bool

func (g GraphLiteralFixture) Vertices(f VertexStep) {
	vl := []Vertex{"foo", "bar", "baz", "isolate"}
	for _, v := range vl {
		if f(v) {
			return
		}
	}
}

func (g GraphLiteralFixture) EachEdge(f EdgeStep) {
	var el []Edge
	if g {
		el = []Edge{
			NewEdge("foo", "bar"),
			NewEdge("bar", "baz"),
		}
	} else {
		el = []Edge{
			NewEdge("bar", "foo"),
			NewEdge("baz", "bar"),
		}
	}

	for _, e := range el {
		if f(e) {
			return
		}
	}
}

func (g GraphLiteralFixture) EachArc(f ArcStep) {
	var al []Arc
	if g {
		al = []Arc{
			NewArc("foo", "bar"),
			NewArc("bar", "baz"),
		}
	} else {
		al = []Arc{
			NewArc("bar", "foo"),
			NewArc("baz", "bar"),
		}
	}

	for _, e := range al {
		if f(e) {
			return
		}
	}
}

func (g GraphLiteralFixture) IncidentTo(v Vertex, f EdgeStep) {
	if g {
		switch v {
		case "foo":
			f(NewEdge("foo", "bar"))
		case "bar":
			terminate := f(NewEdge("foo", "bar"))
			if !terminate {
				f(NewEdge("bar", "baz"))
			}
		case "baz":
			f(NewEdge("bar", "baz"))
		default:
		}
	} else {
		switch v {
		case "foo":
			f(NewEdge("bar", "foo"))
		case "bar":
			terminate := f(NewEdge("bar", "foo"))
			if !terminate {
				f(NewEdge("baz", "bar"))
			}
		case "baz":
			f(NewEdge("baz", "bar"))
		default:
		}
	}
}

func (g GraphLiteralFixture) ArcsFrom(v Vertex, f ArcStep) {
	if g {
		switch v {
		case "foo":
			f(NewArc("foo", "bar"))
		case "bar":
			f(NewArc("bar", "baz"))
		default:
		}
	} else {
		switch v {
		case "bar":
			f(NewArc("bar", "foo"))
		case "baz":
			f(NewArc("baz", "bar"))
		default:
		}
	}
}

func (g GraphLiteralFixture) ArcsTo(v Vertex, f ArcStep) {
	if g {
		switch v {
		case "bar":
			f(NewArc("foo", "bar"))
		case "baz":
			f(NewArc("bar", "baz"))
		default:
		}
	} else {
		switch v {
		case "foo":
			f(NewArc("bar", "foo"))
		case "bar":
			f(NewArc("baz", "bar"))
		default:
		}
	}
}

func (g GraphLiteralFixture) EachPredecessorOf(v Vertex, f VertexStep) {
	if g {
		switch v {
		case "bar":
			f("foo")
		case "baz":
			f("bar")
		default:
		}
	} else {
		switch v {
		case "foo":
			f("bar")
		case "bar":
			f("baz")
		default:
		}
	}
}
func (g GraphLiteralFixture) EachSuccessorOf(v Vertex, f VertexStep) {
	if g {
		switch v {
		case "foo":
			f("bar")
		case "bar":
			f("baz")
		default:
		}
	} else {
		switch v {
		case "bar":
			f("foo")
		case "baz":
			f("bar")
		default:
		}
	}
}

func (g GraphLiteralFixture) EachAdjacentTo(v Vertex, f VertexStep) {
	switch v {
	case "foo":
		f("bar")
	case "bar":
		terminate := f("foo")
		if !terminate {
			f("baz")
		}
	case "baz":
		f("bar")
	default:
	}
}

func (g GraphLiteralFixture) HasVertex(v Vertex) bool {
	switch v {
	case "foo", "bar", "baz", "isolate":
		return true
	default:
		return false
	}
}

func (g GraphLiteralFixture) InDegreeOf(v Vertex) (degree int, exists bool) {
	if g {
		switch v {
		case "foo":
			return 0, true
		case "bar":
			return 1, true
		case "baz":
			return 1, true
		case "isolate":
			return 0, true
		default:
			return 0, false
		}
	} else {
		switch v {
		case "foo":
			return 1, true
		case "bar":
			return 1, true
		case "baz":
			return 0, true
		case "isolate":
			return 0, true
		default:
			return 0, false
		}
	}
}

func (g GraphLiteralFixture) OutDegreeOf(v Vertex) (degree int, exists bool) {
	if g {
		switch v {
		case "foo":
			return 1, true
		case "bar":
			return 1, true
		case "baz":
			return 0, true
		case "isolate":
			return 0, true
		default:
			return 0, false
		}

	} else {
		switch v {
		case "foo":
			return 0, true
		case "bar":
			return 1, true
		case "baz":
			return 1, true
		case "isolate":
			return 0, true
		default:
			return 0, false
		}
	}
}

func (g GraphLiteralFixture) DegreeOf(v Vertex) (degree int, exists bool) {
	switch v {
	case "foo":
		return 1, true
	case "bar":
		return 2, true
	case "baz":
		return 1, true
	case "isolate":
		return 0, true
	default:
		return 0, false
	}
}

func (g GraphLiteralFixture) HasEdge(e Edge) bool {
	u, v := e.Both()

	switch u {
	case "foo":
		return v == "bar"
	case "bar":
		return v == "baz" || v == "foo"
	case "baz":
		return v == "bar"
	default:
		return false
	}
}

func (g GraphLiteralFixture) HasArc(a Arc) bool {
	u, v := a.Both()

	if g {
		switch u {
		case "foo":
			return v == "bar"
		case "bar":
			return v == "baz"
		default:
			return false
		}
	} else {
		switch u {
		case "bar":
			return v == "foo"
		case "baz":
			return v == "bar"
		default:
			return false
		}
	}
}

func (g GraphLiteralFixture) Density() float64 {
	return 2 / 12 // 2 edges of maximum 12 in a 4-vertex digraph
}

func (g GraphLiteralFixture) Transpose() Digraph {
	return GraphLiteralFixture(!g)
}

func (g GraphLiteralFixture) Size() int {
	return 2
}

func (g GraphLiteralFixture) Order() int {
	return 4
}
