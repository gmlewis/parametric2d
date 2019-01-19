// Package parametric2d defines 2D line and cubic curve segments that can be
// interpolated.  It also calculates tangents, normals, and offsets
// to the segments in order to create 3D bevels.
package parametric2d

import (
	"github.com/gmlewis/go-poly2tri"
	"github.com/gmlewis/go3d/float64/vec2"
	"github.com/gmlewis/go3d/float64/vec3"
)

// T represents a parametric 2D segment.
type T interface {
	// BBox returns the bounds of the segment.
	BBox() vec2.Rect
	// At interpolates along the segment and returns the 2D point.
	At(t float64) vec2.T
	// Tangent interpolates along the segment and returns the un-normalized tangent at that point.
	Tangent(t float64) vec2.T
	// NTangent interpolates along the segment and returns the normalized tangent at that point.
	NTangent(t float64) vec2.T
	// Normal interpolates along the segment and returns the un-normalized normal at that point.
	Normal(t float64) vec2.T
	// NNormal interpolates along the segment and returns the normalized normal at that point.
	NNormal(t float64) vec2.T
	// Wall returns a triangularized vertical extrusion of the 2D segment of the given height.
	// It subdivides the segment into as many vertical slices (making 2 triangles out of each slice)
	// such that the maximum angle between adjacent angles is 'maxDegrees'.
	// flipNormals determines if the normals are flipped from their default orientation.
	Wall(height, maxDegrees float64, flipNormals bool) ([]Triangle3D, poly2tri.PointArray)
	// Bevel returns a triangularized angled extrusion of the 2D segment starting at the given height.
	// It subdivides the segment in the same manner as Wall().
	// 'offset' specifies the horizontal distance to offset the original segment.
	// 'deg' specifies the angle (in degrees, where 0 is horizontally flat) to place the bevel.
	// flipNormals determines if the normals are flipped from their default orientation.
	// prevNN is the previous segment's normalized normal at its t=1 endpoint.
	// nextNN is the next segment's normalized normal at its t=0 endpoint..
	Bevel(height, offset, deg, maxDegrees float64, flipNormals bool, prevNN, nextNN *vec2.T) ([]Triangle3D, poly2tri.PointArray)
	// IsLine returns true if this segment is a simple line segment
	IsLine() bool
}

// Triangle3D represents a 3D triangle.
type Triangle3D []vec3.T

// TriangleWriter writes a triangle to some output.
type TriangleWriter interface {
	WriteTriangle3D(t Triangle3D) error
}
