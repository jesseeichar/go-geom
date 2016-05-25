package xygraph

type EdgeSetIntersector interface {
	computeIntersections(edges []Edge, si SegmentIntersector, testAllSegments bool)
	computeIntersectionsForEdges(edges0, edges1 []Edge, si SegmentIntersector)
}