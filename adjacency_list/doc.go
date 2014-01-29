/*
Adjacency lists are a relatively simple graph representation. They maintain
a list of vertices, storing information about edge membership relative to
those vertices. This makes vertex-centric operations generally more
efficient, and edge-centric operations generally less efficient, as edges
are represented implicitly. It also makes them inappropriate for more
complex graph types, such as multigraphs or those with edges having multiple
properties.

gogl's adjacency lists are space-efficient; in a directed graph, the memory
cost for the entire graph G is proportional to V + E; in an undirected graph,
it is V + 2E.

*/
package adjacency_list
