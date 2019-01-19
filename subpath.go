package parametric2d

import (
	"fmt"
	"math"

	"github.com/gmlewis/go-poly2tri"
	"github.com/gmlewis/go3d/float64/vec2"
)

// SubPath represents a 2D collection of parametric segments.
type SubPath struct {
	Segments    []T
	FlipNormals bool
	// IsOuter determines if this is the enclosing outer subpath.
	IsOuter  bool
	BevelPts poly2tri.PointArray
	BevelZ   float64
	FloorPts poly2tri.PointArray
	FloorZ   float64
}

// BBox returns the minimum bounding box of the SubPath.
func (s *SubPath) BBox() vec2.Rect {
	if len(s.Segments) == 0 {
		return vec2.Rect{}
	}
	bbox := s.Segments[0].BBox()
	for _, sp := range s.Segments[1:] {
		v := sp.BBox()
		bbox = vec2.Joined(&bbox, &v)
	}
	return bbox
}

// Wall extrudes a subpath into a 3D wall. `maxDegrees` determines the smoothness
// of the wall along the subpath.
func (s *SubPath) Wall(height, maxDegrees float64) []Triangle3D {
	s.FloorZ = 0
	r := make([]Triangle3D, 0, 100)
	for _, seg := range s.Segments {
		w, floorPts := seg.Wall(height, maxDegrees, s.FlipNormals)
		r = append(r, w...)
		s.FloorPts = append(s.FloorPts, floorPts...)
		// fmt.Printf("GML: subPath floorPts=%#v\n", floorPts)
	}
	return r
}

// Bevel returns a 3D beveled object based on the provided subpath.
func (s *SubPath) Bevel(height, offset, deg, maxDegrees float64) []Triangle3D {
	s.BevelZ = height + offset*math.Tan(deg*math.Pi/180.0)
	r := []Triangle3D{}
	fmt.Printf("\nGML: ENTER Subpath.Bevel: #Segments=%v", len(s.Segments))
	for i, seg := range s.Segments {
		j := (i + len(s.Segments) - 1) % len(s.Segments)
		k := (i + 1) % len(s.Segments)
		prevNN := s.Segments[j].NNormal(1)
		nextNN := s.Segments[k].NNormal(0)
		if s.FlipNormals {
			prevNN[0], prevNN[1], nextNN[0], nextNN[1] = -prevNN[0], -prevNN[1], -nextNN[0], -nextNN[1]
		}
		{ //GML... debug
			n0 := s.Segments[i].NNormal(0)
			n1 := s.Segments[i].NNormal(1)
			fmt.Printf("GML: j=%v, prevNN=%v, i=%v, n0=%v, n1=%v, k=%v, nextNN=%v\n", j, prevNN, i, n0, n1, k, nextNN)
		}
		b, bevelPts := seg.Bevel(height, offset, deg, maxDegrees, s.FlipNormals, &prevNN, &nextNN)
		r = append(r, b...)
		s.BevelPts = append(s.BevelPts, bevelPts...)
		// fmt.Printf("GML: subPath bevelPts={")
		// for _, v := range bevelPts {
		// 	fmt.Printf("%v ", v)
		// }
		// fmt.Printf("}\n")
	}
	return r
}

// AutoFlipNormals analyzes the subpath to determine if its normals
// are currently facing the correct direction (outward), and if not,
// sets the 'FlipNormals' flag on the subpath.
func (s *SubPath) AutoFlipNormals() {
	bbox := s.BBox()
	flippedBBox := bbox
	for _, sp := range s.Segments {
		p0 := sp.At(0)
		p1 := sp.At(1)
		n0 := sp.NNormal(0)
		n1 := sp.NNormal(1)
		newP0 := vec2.Add(&p0, &n0)
		newP1 := vec2.Add(&p1, &n1)
		b := vec2.NewRect(&newP0, &newP1)
		flippedBBox = vec2.Joined(&flippedBBox, &b)
	}
	if bbox.Contains(&flippedBBox) {
		return
	}
	// fmt.Printf("setting flip normals=true, bbox=%v, flippedBBox=%v\n", bbox, flippedBBox)
	s.FlipNormals = true
}
