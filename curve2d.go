package parametric2d

import (
	"fmt"
	"math"

	"github.com/gmlewis/go-poly2tri"
	"github.com/gmlewis/go3d/float64/bezier2"
	"github.com/gmlewis/go3d/float64/vec2"
	"github.com/gmlewis/go3d/float64/vec3"
)

// Curve represents a 2D Bezier curve and implements interface T.
type Curve struct {
	spline bezier2.T
	bbox   vec2.Rect
}

// NewCurve returns a new 2D Bezier curve from four points.
func NewCurve(p0, p1, p2, p3 vec2.T) Curve {
	ll := vec2.Min(&p0, &p3)
	ur := vec2.Max(&p0, &p3)
	// Evaluate the bezier spline at t=0.25, 0.5, 0.75
	// to get accurate bounds for the curve.
	s := bezier2.T{P0: p0, P1: p1, P2: p2, P3: p3}
	if p0 == p1 {
		fmt.Printf("Warning! cubic bezier p0==p1==%v.  Setting p1=p2=%v\n", p0, p2)
		s.P1 = p2
	}
	if p2 == p3 {
		fmt.Printf("Warning! cubic bezier p2==p3==%v.  Setting p2=p1=%v\n", p3, p1)
		s.P2 = p1
	}
	for _, t := range []float64{0.25, 0.5, 0.75} {
		v := s.Point(t)
		ll = vec2.Min(&ll, &v)
		ur = vec2.Max(&ur, &v)
	}
	return Curve{spline: s, bbox: vec2.Rect{Min: ll, Max: ur}}
}

// BBox returns the minimum bounding box of the Curve.
func (s Curve) BBox() vec2.Rect {
	return s.bbox
}

// At returns the point on the Curve at the given position (0 <= t <= 1).
func (s Curve) At(t float64) vec2.T {
	return s.spline.Point(t)
}

// Tangent returns the tangent to the Curve
// at the given position (0 <= t <= 1).
func (s Curve) Tangent(t float64) vec2.T {
	return s.spline.Tangent(t)
}

// NTangent returns the normalized tangent to the Curve
// at the given position (0 <= t <= 1).
func (s Curve) NTangent(t float64) vec2.T {
	v := s.Tangent(t)
	return *v.Normalize()
}

// Normal returns the normal to the Curve
// at the given position (0 <= t <= 1).
func (s Curve) Normal(t float64) vec2.T {
	v := s.Tangent(t)
	return *v.Rotate90DegLeft()
}

// NNormal returns the normalized normal to the Curve
// at the given position (0 <= t <= 1).
func (s Curve) NNormal(t float64) vec2.T {
	v := s.Normal(t)
	return *v.Normalize()
}

// Subdivide returns the parametric 't' values along the curve
// such that the tangent between two points never exceeds `maxDegrees`.
func (s Curve) Subdivide(maxDegrees float64) []float64 {
	// Start with 3 points, and subdivide as necessary:
	ts := []float64{0, 0.5, 1}
	tangents := []vec2.T{s.NTangent(0), s.NTangent(0.5), s.NTangent(1)}
	maxRadians := math.Abs(maxDegrees * math.Pi / 180.0)
	i := 0
	for i < len(ts)-1 {
		if ts[i+1]-ts[i] < 1e-2 {
			fmt.Printf("Stopping subdivision: ts[%v]=%v, ts[%v]=%v\n", i+1, ts[i+1], i, ts[i])
			i++
			continue
		}
		angle := math.Abs(vec2.Angle(&tangents[i], &tangents[i+1]))
		if angle > maxRadians { // Subdivide
			m := 0.5 * (ts[i] + ts[i+1])
			tan := s.NTangent(m)
			// Insert the new t and tangent into the slice.
			ts = append(ts, 0) // make room
			copy(ts[i+2:], ts[i+1:])
			ts[i+1] = m
			tangents = append(tangents, vec2.T{0, 0}) // make room
			copy(tangents[i+2:], tangents[i+1:])
			tangents[i+1] = tan
		} else {
			i++
		}
	}
	return ts
}

// Wall extrudes a curve into a 3D wall. `maxDegrees` determines the smoothness
// of the wall along the curve.
func (s Curve) Wall(height, maxDegrees float64, flipNormals bool) ([]Triangle3D, poly2tri.PointArray) {
	ts := s.Subdivide(maxDegrees)
	num := len(ts)
	if num <= 0 {
		return []Triangle3D{}, poly2tri.PointArray{}
	}
	v := make([]Triangle3D, 0, 2*num)
	floorPts := make(poly2tri.PointArray, 0, num-1)
	for i := 0; i < num-1; i++ {
		p0 := s.At(ts[i])
		p1 := s.At(ts[i+1])
		t := Triangle3D{
			vec3.T{p0[0], p0[1], 0},
			vec3.T{p1[0], p1[1], height},
			vec3.T{p0[0], p0[1], height},
		}
		if flipNormals {
			t[1], t[2] = t[2], t[1]
		}
		v = append(v, t)
		// if i == 0 {
		// 	floorPts = append(floorPts, poly2tri.NewPoint(p0[0], p0[1]))
		// }
		t = Triangle3D{
			vec3.T{p0[0], p0[1], 0},
			vec3.T{p1[0], p1[1], 0},
			vec3.T{p1[0], p1[1], height},
		}
		if flipNormals {
			t[1], t[2] = t[2], t[1]
		}
		v = append(v, t)
		floorPts = append(floorPts, poly2tri.NewPoint(p1[0], p1[1]))
	}
	return v, floorPts
}

// Bevel returns a 3D beveled object based on the provided curve.
func (s Curve) Bevel(height, offset, deg, maxDegrees float64, flipNormals bool, prevNN, nextNN *vec2.T) ([]Triangle3D, poly2tri.PointArray) {
	ts := s.Subdivide(maxDegrees)
	num := len(ts)
	if num <= 0 {
		return []Triangle3D{}, poly2tri.PointArray{}
	}
	h := offset * math.Tan(deg*math.Pi/180.0)
	v := make([]Triangle3D, 0, 2*num)
	bevelPts := make(poly2tri.PointArray, 0, num-1)
	for i := 0; i < num-1; i++ {
		p0 := s.At(ts[i])
		p1 := s.At(ts[i+1])
		n0 := s.NNormal(ts[i])
		n1 := s.NNormal(ts[i+1])
		n0f := offset
		n1f := offset
		if flipNormals {
			n0[0], n0[1], n1[0], n1[1] = -n0[0], -n0[1], -n1[0], -n1[1]
		}
		if i == 0 {
			angle0 := vec2.Angle(prevNN, &n0)
			if *prevNN != n0 {
				// Adjusting starting triangle n0f=1.0000532533019526, prevNN=[0.3559858534155444 -0.9344913440840459], n0=[0.37519651438861046 -0.9269452926632926], angle0=0.020639949379423816
				// Created regular start-of-curve triangle:
				// [[181.08017999999998 -499.24048 4] [182.489264893852 -499.72347812262666 5] [181.44581012281427 -500.17129744866037 5]]
				n0f = offset / math.Cos(0.5*angle0)
				fmt.Printf("Adjusting starting triangle n0f=%v, prevNN=%v, n0=%v, angle0=%v\n", n0f, *prevNN, n0, angle0)
				n0.Add(prevNN)
				n0.Normalize()
			} else {
				fmt.Printf("Warning: prevNN==n0==%v, i=%v, p0=%v\n", n0, i, p0)
			}
		}
		if i+1 == num-1 {
			angle1 := vec2.Angle(&n1, nextNN)
			if n1 != *nextNN {
				n1f = offset / math.Cos(0.5*angle1)
				// fmt.Printf("n1f=%v, n1=%v, nextNN=%v, angle1=%v\n", n1f, n1, *nextNN, angle1)
				n1.Add(nextNN)
				n1.Normalize()
			} else {
				fmt.Printf("Warning: n1==nextNN==%v, i=%v, p1=%v\n", n1, i, p1)
			}
		}
		p2 := n0.Scale(n0f).Add(&p0)
		p3 := n1.Scale(n1f).Add(&p1)
		if segmentsIntersect(&p0, p2, &p1, p3) { // p0-p2 intersects p1-p3 - delete p3 and add new triangle
			fmt.Printf("Detected intersection: i=%v, num=%v\n", i, num)
			if i == 0 {
				// Do not add a bevelPts for this case because p2 has already been accounted for by last segment
				t := Triangle3D{
					vec3.T{p0[0], p0[1], height},
					vec3.T{p1[0], p1[1], height},
					vec3.T{p2[0], p2[1], height + h},
				}
				fmt.Printf("Created start-of-curve triangle: %v\n", t)
				if flipNormals {
					t[1], t[2] = t[2], t[1]
				}
				v = append(v, t)
			} else if i+1 == num-1 { // End of the curve - need to add to the bevelPts
				t := Triangle3D{
					vec3.T{p0[0], p0[1], height},
					vec3.T{p1[0], p1[1], height},
					vec3.T{p3[0], p3[1], height + h},
				}
				fmt.Printf("Created end-of-curve triangle: %v\n", t)
				// Created end-of-curve triangle:
				// Intersection: [182.07626125,-498.81274874999997]-[182.489264893852,-499.72347812262666]
				//             X [183.09580999999997,-498.32642]-[182.09580999999997,-499.9444329390584]
				// [[182.07626125 -498.81274874999997 4]
				// [183.09580999999997 -498.32642 4]
				// [182.09580999999997 -499.9444329390584 5]]
				if flipNormals {
					t[1], t[2] = t[2], t[1]
				}
				v = append(v, t)
				bevelPts = append(bevelPts, poly2tri.NewPoint(p3[0], p3[1]))
			} else { // adjust to previous point.
				panic("curve2d.go: unhandled problem")
			}
		} else {
			t := Triangle3D{
				vec3.T{p0[0], p0[1], height},
				vec3.T{p3[0], p3[1], height + h},
				vec3.T{p2[0], p2[1], height + h},
			}
			if i == 0 {
				fmt.Printf("Created regular start-of-curve triangle:\n%v\n", t)
			}
			if flipNormals {
				t[1], t[2] = t[2], t[1]
			}
			v = append(v, t)
			t = Triangle3D{
				vec3.T{p0[0], p0[1], height},
				vec3.T{p1[0], p1[1], height},
				vec3.T{p3[0], p3[1], height + h},
			}
			if flipNormals {
				t[1], t[2] = t[2], t[1]
			}
			v = append(v, t)
			bevelPts = append(bevelPts, poly2tri.NewPoint(p3[0], p3[1]))
		}
	}
	fmt.Printf("Bevel: #subdivisions=%v: #v=%v, #bevelPts=%v\n", num, len(v), len(bevelPts))
	return v, bevelPts
}

// IsLine is false for type Curve.
func (s Curve) IsLine() bool {
	return false
}

func segmentsIntersect(a, b, c, d *vec2.T) bool {
	div := (a[0]-b[0])*(c[1]-d[1]) - (a[1]-b[1])*(c[0]-d[0])
	if div == 0 {
		return false
	}
	x := ((a[0]*b[1]-a[1]*b[0])*(c[0]-d[0]) - (a[0]-b[0])*(c[0]*d[1]-c[1]*d[0])) / div
	y := ((a[0]*b[1]-a[1]*b[0])*(c[1]-d[1]) - (a[1]-b[1])*(c[0]*d[1]-c[1]*d[0])) / div
	if a[0] >= b[0] {
		if !(b[0] <= x && x <= a[0]) {
			return false
		}
	} else {
		if !(a[0] <= x && x <= b[0]) {
			return false
		}
	}
	if a[1] >= b[1] {
		if !(b[1] <= y && y <= a[1]) {
			return false
		}
	} else {
		if !(a[1] <= y && y <= b[1]) {
			return false
		}
	}
	if c[0] >= d[0] {
		if !(d[0] <= x && x <= c[0]) {
			return false
		}
	} else {
		if !(c[0] <= x && x <= d[0]) {
			return false
		}
	}
	if c[1] >= d[1] {
		if !(d[1] <= y && y <= c[1]) {
			return false
		}
	} else {
		if !(c[1] <= y && y <= d[1]) {
			return false
		}
	}
	fmt.Printf("\nIntersection: [%v,%v]-[%v,%v] X [%v,%v]-[%v,%v]\n", a[0], a[1], b[0], b[1], c[0], c[1], d[0], d[1])
	return true
}
