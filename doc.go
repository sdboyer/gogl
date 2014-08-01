/*
gogl provides a framework for representing and working with graphs.

gogl is divided chiefly into two parts: datastructures and algorithms.
These components are loosely coupled; the algorithms rely strictly on
gogl's various exported Graph interfaces to do their work.

Idiomatic Go emphasizes small (1-2 method), purpose-oriented interfaces.
gogl takes this to an extreme; the main graph package exports more than
30 interfaces. While this has great benefit in terms of correctness
and flexibility, it is daunting - especially if you're looking at the
godoc, where everthing is an alpha-sorted mess.

In the source, however, human-friendly order is intentionally maintained.
To learn gogl by reading the code (recommended!), you can learn the
concepts incrementally by reading these files, in order and from top to
bottom:

- vertex.go
- edge.go
- graph.go

*/
package gogl
