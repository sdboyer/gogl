package gogl

var _ = SetUpSimpleGraphTests(NewWeightedDirected(), true)
var _ = SetUpSimpleGraphTests(NewWeightedUndirected(), false)
