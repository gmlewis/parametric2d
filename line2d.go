package parametric2d

import (
	"math"

	"github.com/gmlewis/go-poly2tri"
	"github.com/gmlewis/go3d/float64/vec2"
	"github.com/gmlewis/go3d/float64/vec3"
)

// Line represents a straight line segment and implements interface T.
type Line struct {
	p0, p1 vec2.T
	bbox   vec2.Rect
}

// NewLine returns a new 2D Line from two points.
func NewLine(p0, p1 vec2.T) Line {
	ll := vec2.Min(&p0, &p1)
	ur := vec2.Max(&p0, &p1)
	return Line{p0: p0, p1: p1, bbox: vec2.Rect{Min: ll, Max: ur}}
}

// BBox returns the minimum bounding box of the Line.
func (s Line) BBox() vec2.Rect {
	return s.bbox
}

// At returns the point on the Line at the given position (0 <= t <= 1).
func (s Line) At(t float64) vec2.T {
	return vec2.Interpolate(&s.p0, &s.p1, t)
}

// Tangent returns the tangent to the Line
// at the given position (0 <= t <= 1).
func (s Line) Tangent(t float64) vec2.T {
	return vec2.Sub(&s.p1, &s.p0)
}

// NTangent returns the normalized tangent to the Line
// at the given position (0 <= t <= 1).
func (s Line) NTangent(t float64) vec2.T {
	v := s.Tangent(t)
	return *v.Normalize()
}

// Normal returns the normal to the Line
// at the given position (0 <= t <= 1).
func (s Line) Normal(t float64) vec2.T {
	v := s.Tangent(t)
	return *v.Rotate90DegLeft()
}

// NNormal returns the normalized normal to the Line
// at the given position (0 <= t <= 1).
func (s Line) NNormal(t float64) vec2.T {
	v := s.Normal(t)
	return *v.Normalize()
}

// Wall extrudes a line into a 3D wall. `maxDegrees` is ignored.
func (s Line) Wall(height, maxDegrees float64, flipNormals bool) ([]Triangle3D, poly2tri.PointArray) {
	p0 := s.At(0)
	p1 := s.At(1)
	t0 := Triangle3D{
		vec3.T{p0[0], p0[1], 0},
		vec3.T{p1[0], p1[1], height},
		vec3.T{p0[0], p0[1], height},
	}
	if flipNormals {
		t0[1], t0[2] = t0[2], t0[1]
	}
	t1 := Triangle3D{
		vec3.T{p0[0], p0[1], 0},
		vec3.T{p1[0], p1[1], 0},
		vec3.T{p1[0], p1[1], height},
	}
	if flipNormals {
		t1[1], t1[2] = t1[2], t1[1]
	}
	return []Triangle3D{t0, t1}, poly2tri.PointArray{poly2tri.NewPoint(p1[0], p1[1])}
}

// Bevel returns a 3D beveled object based on the provided Line.
func (s Line) Bevel(height, offset, deg, maxDegrees float64, flipNormals bool, prevNN, nextNN *vec2.T) ([]Triangle3D, poly2tri.PointArray) {
	h := offset * math.Tan(deg*math.Pi/180.0)
	p0 := s.At(0)
	p1 := s.At(1)
	n0 := s.NNormal(0)
	n1 := s.NNormal(1)
	if flipNormals {
		n0[0], n0[1], n1[0], n1[1] = -n0[0], -n0[1], -n1[0], -n1[1]
	}
	angle0 := vec2.Angle(prevNN, &n0)
	newLength0 := offset / math.Cos(0.5*angle0)
	angle1 := vec2.Angle(&n1, nextNN)
	newLength1 := offset / math.Cos(0.5*angle1)
	prevNN.Add(&n0)
	prevNN.Normalize()
	nextNN.Add(&n1)
	nextNN.Normalize()
	p2 := prevNN.Scale(newLength0).Add(&p0)
	p3 := nextNN.Scale(newLength1).Add(&p1)
	t0 := Triangle3D{
		vec3.T{p0[0], p0[1], height},
		vec3.T{p3[0], p3[1], height + h},
		vec3.T{p2[0], p2[1], height + h},
	}
	if flipNormals {
		t0[1], t0[2] = t0[2], t0[1]
	}
	t1 := Triangle3D{
		vec3.T{p0[0], p0[1], height},
		vec3.T{p1[0], p1[1], height},
		vec3.T{p3[0], p3[1], height + h},
	}
	if flipNormals {
		t1[1], t1[2] = t1[2], t1[1]
	}
	return []Triangle3D{t0, t1}, poly2tri.PointArray{poly2tri.NewPoint(p3[0], p3[1])}
}

// IsLine is true for type Line.
func (s Line) IsLine() bool {
	return true
}
