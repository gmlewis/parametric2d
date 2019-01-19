package parametric2d

import (
	"fmt"
	"sort"

	"github.com/gmlewis/go-poly2tri"
	"github.com/gmlewis/go3d/float64/vec2"
	"github.com/gmlewis/go3d/float64/vec3"
)

// Path represents a 2D collection of SubPaths.
type Path struct {
	SubPaths []*SubPath
}

// BBox returns the minimum bounding box of the Path.
func (p *Path) BBox() vec2.Rect {
	if len(p.SubPaths) == 0 {
		return vec2.Rect{}
	}
	bbox := p.SubPaths[0].BBox()
	for _, sp := range p.SubPaths[1:] {
		v := sp.BBox()
		bbox = vec2.Joined(&bbox, &v)
	}
	return bbox
}

// Wall extrudes a path into a 3D wall. `maxDegrees` determines the smoothness
// of the wall along the path.
func (p *Path) Wall(height, maxDegrees float64) []Triangle3D {
	if len(p.SubPaths) == 0 {
		return []Triangle3D{}
	}
	bbox := p.SubPaths[0].BBox()
	r := make([]Triangle3D, 0, 100)
	var sc *poly2tri.SweepContext
	for i, sp := range p.SubPaths {
		w := sp.Wall(height, maxDegrees)
		r = append(r, w...)
		if i == 0 {
			fmt.Printf("Initializing %v Floor points, #p.SubPaths=%v\n", len(sp.FloorPts), len(p.SubPaths))
			sc = poly2tri.New(sp.FloorPts)
		} else {
			subBBox := sp.BBox()
			fmt.Printf("Checking if floor bbox %v contains subBBox %v...\n", bbox, subBBox)
			if bbox.Contains(&subBBox) {
				fmt.Printf("Adding %v Floor hole points\n", len(sp.FloorPts))
				sc.AddHole(sp.FloorPts)
			} else {
				fmt.Printf("SubPath not contained by parent... Adding %v new Floor points\n", len(sp.FloorPts))
				r = Triangulate(sc.Triangulate(), r, p.SubPaths[0].FloorZ)
				sc = poly2tri.New(sp.FloorPts)
			}
		}
	}
	r = Triangulate(sc.Triangulate(), r, p.SubPaths[0].FloorZ)
	return r
}

// Bevel returns a 3D beveled object based on the provided path.
func (p *Path) Bevel(height, offset, deg, maxDegrees float64) []Triangle3D {
	if len(p.SubPaths) == 0 {
		return []Triangle3D{}
	}
	bbox := p.SubPaths[0].BBox()
	r := make([]Triangle3D, 0, 100)
	var sc *poly2tri.SweepContext
	for i, sp := range p.SubPaths {
		w := sp.Bevel(height, offset, deg, maxDegrees)
		r = append(r, w...)
		if i == 0 {
			fmt.Printf("Initializing %v Bevel points, #p.SubPaths=%v\n", len(sp.BevelPts), len(p.SubPaths))
			if len(p.SubPaths) == 1 {
				fmt.Printf("\nvar in = PointArray{\n")
				for _, t := range sp.BevelPts {
					fmt.Printf("  NewPoint(%v, %v),\n", t.X, t.Y)
				}
				fmt.Printf("}\n\n")
			}
			sc = poly2tri.New(sp.BevelPts)
		} else {
			subBBox := sp.BBox()
			fmt.Printf("Checking if bevel bbox %v contains subBBox %v...\n", bbox, subBBox)
			if bbox.Contains(&subBBox) {
				fmt.Printf("Adding %v Bevel hole points\n", len(sp.BevelPts))
				sc.AddHole(sp.BevelPts)
			} else {
				fmt.Printf("SubPath not contained by parent... Adding %v new Bevel points\n", len(sp.BevelPts))
				r = Triangulate(sc.Triangulate(), r, p.SubPaths[0].BevelZ)
				sc = poly2tri.New(sp.BevelPts)
			}
		}
	}
	r = Triangulate(sc.Triangulate(), r, p.SubPaths[0].BevelZ)
	return r
}

// Triangulate converts 2D points to 3D triangles and appends them to a slice.
func Triangulate(m poly2tri.TriArray, r []Triangle3D, z float64) []Triangle3D {
	for _, t := range m {
		p0 := vec3.T{t.Point[0].X, t.Point[0].Y, z}
		p1 := vec3.T{t.Point[1].X, t.Point[1].Y, z}
		p2 := vec3.T{t.Point[2].X, t.Point[2].Y, z}
		tri := Triangle3D{p0, p1, p2}
		r = append(r, tri)
	}
	return r
}

// AutoFlipNormals analyzes the first subpath to determine if its normals
// are currently facing the correct direction (outward), and if not,
// sets the 'FlipNormals' flag on each segment of each subpath.
func (p *Path) AutoFlipNormals() {
	if len(p.SubPaths) == 0 {
		return
	}
	// First, sort the subpaths by bounding box square area... Largest first
	sort.Sort(byBBoxArea(p.SubPaths))

	p.SubPaths[0].IsOuter = true
	p.SubPaths[0].AutoFlipNormals()
	if p.SubPaths[0].FlipNormals {
		for _, sp := range p.SubPaths[1:] {
			sp.FlipNormals = true
		}
	}
}

type byBBoxArea []*SubPath // slice of *SubPath

func (x byBBoxArea) Len() int      { return len(x) }
func (x byBBoxArea) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byBBoxArea) Less(i, j int) bool {
	iBBox := x[i].BBox()
	jBBox := x[j].BBox()
	return jBBox.Area() < iBBox.Area()
}
