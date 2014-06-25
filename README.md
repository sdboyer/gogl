# gogl

[![Build Status](https://travis-ci.org/sdboyer/gogl.png?branch=master)](https://travis-ci.org/sdboyer/gogl)
[![Coverage Status](https://coveralls.io/repos/sdboyer/gogl/badge.png?branch=master)](https://coveralls.io/r/sdboyer/gogl?branch=master)

gogl is a graph library in Go. Its goal is to provide simple, unifying interfaces and implementations of graph algorithms and datastructures that can scale from use cases for small graphs up to very large (read: distributed) graphs.

gogl is based on the premise that graph computation problems can be [decomplected](http://www.infoq.com/presentations/Simple-Made-Easy) by focusing, first and foremost, on an API modeled on nothing more (or less) than the intrinsic beauty of graph theory - in essence, an isomorphism. Certainly a lofty goal, and possibly a silly premise.

There's still a lot to do - gogl is still firming up significant aspects of how its API works.

## Principles

These principles guide the way that gogl works:

Graph systems are often big, complicated affairs. gogl is trying to not be that. These are the operant principles:

1. Attain simplicity by expressing the intrinsic nature of graphs through idiomatic Go interface. Other stuff is noise.
1. Be as fast as design constraints and known-best algorithms allow.
1. Expect others to implement their own algorithms using gogl's graph datastructures, and gogl's algorithms with their graph implementations.
1. Build in layers, orient towards transforms, and remain generally functional in style.
1. Don't try to be a graph database (though being the basis for one is totally cool).
1. Be unopinionated about vertices, and minimally opinionated about edges.
1. Utilize [commonly accepted graph terminology](http://en.wikipedia.org/wiki/Glossary_of_graph_theory) where possible.

The [godoc](https://godoc.org/github.com/sdboyer/gogl) contains a lot more discussion about gogl's philosophy. But godoc is a reference, not a tutorial.

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
    // Creates a mutable directed graph
    g := gogl.NewDirected()
    // Adds two edges (and implicitly, their vertices)
    g.AddEdge(gogl.BaseEdge{"foo", "bar"}, gogl.BaseEdge{"bar", "baz"}))

    // Topologically sort your graph from a given start point, producing a slice of vertices
    var tsl []gogl.Vertex
    tsl, _ = dfs.Toposort(g, "foo") // second return is an error
    fmt.Println(tsl) // prints "[baz bar foo]"
}
```

TODO lots and lots and lots more!
