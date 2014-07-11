package gogl

/* Vertex structures */

// A Vertex in gogl is a value of empty interface type.
//
// This is largely comensurate with graph theory, as vertices are, in the general
// case, known to be unique within the graph context, but the particular property
// that makes them unique has no bearing on the graph's behavior. In Go-speak, that
// translates pretty nicely to interface{}.
type Vertex interface{}

