package xygraph

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/location"
	"github.com/twpayne/go-geom/xy/orientation"
)

func linkResultDirectedEdges(nodes []*Node) {
	for _, node := range nodes {
		node.edges.(DirectedEdgeStar).linkResultDirectedEdges()
	}
}

// PlanarGraph contains nodes and edges corresponding to the nodes and line segments of
// a Geometry. Each node and edge in the graph is labeled with its topological location
// relative to the source geometry.
//
// Note that there is no requirement that points of self-intersection be a vertex.
// Thus to obtain a correct topology graph, Geometrys must be
// self-noded before constructing their graphs.
//
// Two fundamental operations are supported by topology graphs:
// * Computing the intersections between all the edges and nodes of a single graph
// * Computing the intersections between the edges and nodes of two different graphs
type PlanarGraph struct {
	edges       []*Edge
	nodes       *NodeMap
	edgeEndList []EdgeEnd
}

func NewPlanarGraph(nodeFactory NodeFactory) *PlanarGraph {
	return &PlanarGraph{
		nodes: &NodeMap{nodeFactory: nodeFactory},
	}
}

func (pg *PlanarGraph) isBoundaryNode(geomIndex int, coord geom.Coord) {
	node, has := pg.nodes.find(coord)
	if !has {
		return false
	}
	label := node.label
	if label != nil && label[geomIndex][ON] == location.Boundary {
		return true
	}
	return false
}

func (pg *PlanarGraph) insertEdge(e *Edge) {
	pg.edges = append(pg.edges, e)
}
func (pg *PlanarGraph) addEdgeEnd(edgeEnd EdgeEnd) {
	pg.nodes.addEdgeEnd(edgeEnd)
	pg.edgeEndList = append(pg.edgeEndList, edgeEnd)
}
func (pg *PlanarGraph) addNode(n *Node) *Node {
	return pg.nodes.addNode(n)
}
func (pg *PlanarGraph) addCoord(c geom.Coord) *Node {
	return pg.nodes.addCoordNode(c)
}
func (pg *PlanarGraph) find(c geom.Coord) {
	pg.nodes.find(c)
}

// addEdges add a set of edges to the graph.  For each edge two DirectedEdges
// will be created.  DirectedEdges are NOT linked by this method.
func (pg *PlanarGraph) addEdges(edgesToAdd []*Edge) {
	// create all the nodes for the edges
	for _, e := range edgesToAdd {
		pg.edges = append(pg.edges, e)

		de1 := NewDirectedEdge(e, true)
		de2 := NewDirectedEdge(e, false)
		de1.sym = de2
		de2.sym = de1

		pg.addEdgeEnd(de1)
		pg.addEdgeEnd(de2)
	}
}

// linkResultDirectedEdges links the DirectedEdges at the nodes of the graph.
// This allows clients to link only a subset of nodes in the graph, for
// efficiency (because they know that only a subset is of interest).
func (pg *PlanarGraph) linkResultDirectedEdges() {
	pg.nodes.nodeMap.Walk(func(c, n interface{}) {
		node := n.(*Node)
		node.edges.(DirectedEdgeStar).linkResultDirectedEdges()
	})
}

// linkAllDirectedEdges link the DirectedEdges at the nodes of the graph.
// This allows clients to link only a subset of nodes in the graph, for
// efficiency (because they know that only a subset is of interest).
func (pg *PlanarGraph) linkAllDirectedEdges() {
	pg.nodes.nodeMap.Walk(func(c, n interface{}) {
		node := n.(*Node)
		node.edges.(DirectedEdgeStar).linkAllDirectedEdges()
	})
}

// findEdgeEnd Returns the EdgeEnd which has edge e as its base edge
// (MD 18 Feb 2002 - this should return a pair of edges)
func (pg *PlanarGraph) findEdgeEnd(e EdgeEnd) (result EdgeEnd, has bool) {
	for _, ee := range pg.edgeEndList {
		if ee.Edge() == e {
			return ee, true
		}
	}

	return nil, false
}

// findEdge finds the edge whose first two coordinates are p0 and p1
func (pg *PlanarGraph) findEdge(p0, p1 geom.Coord) (edge *Edge, has bool) {
	for _, e := range pg.edges {
		if xy.Equal(p0, 0, e.pts[0], 0) && xy.Equal(p1, 0, e.pts[1], 0) {
			return e, true
		}
	}
	return nil, false
}

// findEdgeInSameDirection finds the edge which starts at p0 and whose first segment is parallel to p1
func (pg *PlanarGraph) findEdgeInSameDirection(p0, p1 geom.Coord) (edge *Edge, has bool) {
	for _, e := range pg.edges {
		if pg.matchInSameDirection(p0, p1, e.pts[0], e.pts[1]) {
			return e, true
		}

		if pg.matchInSameDirection(p0, p1, e.pts[len(e.pts)-1], e.pts[len(e.pts)-2]) {
			return e, true
		}
	}
	return nil, false
}

// matchInSameDirection The coordinate pairs match if they define line segments lying in the same direction.
// E.g. the segments are parallel and in the same quadrant (as opposed to parallel and opposite!).
func (pg *PlanarGraph) matchInSameDirection(p0, p1, ep0, ep1 geom.Coord) bool {
	if !xy.Equal(p0, 0, ep0, 0) {
		return false
	}

	if xy.OrientationIndex(p0, p1, ep1) == orientation.Collinear && coordsQuadrant(p0, p1) == coordsQuadrant(ep0, ep1) {
		return true
	}
	return false
}
