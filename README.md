# gogl

[![Build Status](https://travis-ci.org/sdboyer/gogl.png?branch=master)](https://travis-ci.org/sdboyer/gogl)
[![Coverage Status](https://coveralls.io/repos/sdboyer/gogl/badge.png?branch=master)](https://coveralls.io/r/sdboyer/gogl?branch=master)

gogl is a graph library in Go. Its goal is to provide simple, unifying interfaces and implementations of graph algorithms and datastructures that can scale from small graphs to very large graphs. The latter case is, as yet, untested!

gogl is based on the premise that working with graphs can be [decomplected](http://www.infoq.com/presentations/Simple-Made-Easy) by focusing primarily on the natural constraints established in graph theory.

There's still a lot to do - gogl is still firming up significant aspects of how its API works.

## Principles

Graph systems are often big, complicated affairs. gogl tries to be not that. These are the operant principles:

1. Simplicity: fully and correctly modeling graph theoretic concepts in idiomatic Go.
1. Performance: be as fast as design constraints and known-best algorithms allow.
1. Extensibility: expect others to run gogl's graph datastructures through their own algorithms, , and gogl's algorithms with their graph implementations.
1. Functional: orient towards transforms, functors, and streams; achieve other styles through layering.
1. Flexibility: Be unopinionated about vertices, and minimally opinionated about edges.
1. Correctness: Utilize [commonly accepted graph terminology](http://en.wikipedia.org/wiki/Glossary_of_graph_theory) where possible, and adhere to its meaning.

The first and last points are key - names in gogl are carefully chosen, with the hope that they can guide intuition when stricter rules (e.g., the type system) become ambiguous. The [godoc](https://godoc.org/github.com/sdboyer/gogl) generally takes care to detail these subtleties. But godoc is a reference, not a tutorial.

## Quickstart

Getting started with gogl is simple: create a graph object, add your data, and off you go.

```go
package main

import (
	"fmt"
	"github.com/sdboyer/gogl"
	"github.com/sdboyer/gogl/dfs"
)

func main() {
	// gogl uses a builder to specify the kind of graph you want.
	graph := gogl.BuildGraph().
		// Make the graph mutable. default is immutable.
		Mutable().
		// The graph should have directed edges (arcs). default is undirected.
		Directed().
		// The graph's edges are plain - no labels, weights, etc. this is the default.
		Basic().
		// No loops or parallel edges. this is the default.
		SimpleGraph().
		// gogl.AdjacencyList takes the spec selects an adjacency list-based implementation and returns it.
		Create(gogl.AdjacencyList).
		// The builder always returns a Graph; type assert to get access to add/remove methods.
		(gogl.MutableGraph)

	// Adds two basic edges. Of course, this adds the vertices, too.
	graph.AddEdges(gogl.NewEdge("foo", "bar"), gogl.NewEdge("bar", "baz"))

	// gogl's core iteration concept is built on injected functions (VertexLambda or EdgeLambda).
	// Here, a VertexLambda is called once per vertex in the graph; returning 'true' would terminate traversal.
	graph.EachVertex(func(v gogl.Vertex) (terminate bool) {
		fmt.Println(v) // Probably "foo\nbar\nbaz", but ordering is not guaranteed.
		return // returns false, so iteration continues
	})

	// gogl specifies four methods like this on undirected graphs, plus two more on directed graphs.
	// These methods are generally referred to as enumerators; more on them later.

	// If you know you need the full result set, gogl provides functors to collect enumerations
	// into slices. This makes ranging easy.
	var vertices []gogl.Vertex = gogl.CollectVertices(graph)
	for _, v := range vertices {
		fmt.Println(v) // same as with EachVertex().
	}

	// The pattern is the same with edge enumeration.
	graph.EachEdge(func(e gogl.Edge) (terminate bool) {
		fmt.Println(e) // Probably "{foo bar}\n{bar baz}". Again, ordering is not guaranteed.
		return
	})
	for _, e := range gogl.CollectEdges(graph) {
		fmt.Println(e) // same as with EachEdge().
	}

	// gogl's algorithms all rely on these enumerators to do their work. Here, we use
	// a depth-first topological sort algorithm to produce a slice of vertices.
	var tsl []gogl.Vertex
	tsl, err := dfs.Toposort(graph, "foo")
	if err == nil {
		fmt.Println(tsl) // [baz bar foo]
	}
}
```

## Enumerators

TODO - add diagrams indicating what relationship each enumerator touches.
