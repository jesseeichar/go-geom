package lineintersection

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/internal"
)

// Type enumerates the types of intersection two lines can have
type Type int

const (
	// NoIntersection indicates that the lines do not intersect
	NoIntersection Type = iota
	// PointIntersection indicates that the lines intersect at a point
	PointIntersection
	// CollinearIntersection indicates that the lines overlap each other
	CollinearIntersection
)

var labels = [3]string{"NoIntersection", "PointIntersection", "CollinearIntersection"}

func (t Type) String() string {
	return labels[t]
}

// Result the results from LineIntersectsLine function.
// It contains the intersection point(s) and indicates what type of
// intersection there was (or if there was no intersection)
type Result struct {
	isProper         bool
	intersectionType Type
	intersection     []geom.Coord
}

// NewResult create a new result object
func NewResult(intersectionType Type, isProper bool, intersection []geom.Coord) Result {
	return Result{
		intersectionType: intersectionType,
		intersection:     intersection}
}

// HasIntersection returns true if the lines have an intersection
func (i *Result) HasIntersection() bool {
	return i.intersectionType != NoIntersection
}

// Type returns the type of intersection between the two lines
func (i *Result) Type() Type {
	return i.intersectionType
}

// Intersection returns an array of Coords which are the intersection points.
// If the type is PointIntersection then there will only be a single Coordinate (the first coord).
// If the type is CollinearIntersection then there will two Coordinates the start and end points of the line
// that represents the intersection
func (i *Result) Intersection() []geom.Coord {
	return i.intersection
}

// IsProper returns true if the intersection is a "proper" intersection
// The intersection between two line segments is considered proper if
// they intersect in a single point in the interior of both segments
// (e.g. the intersection is a single point and is not equal to any of the endpoints).
func (i *Result) IsProper() bool {
	return i.intersectionType != NoIntersection && i.isProper
}

// IsIntersectionPoint tests if the point is one of the intersection points
func (i *Result) IsIntersectionPoint(point geom.Coord) bool {
	for _, iPt := range i.intersection {
		if internal.Equal(point, 0, iPt, 0) {
			return true
		}
	}
	return false
}
