package parametric2d

import (
	"math"
	"testing"

	"github.com/gmlewis/go3d/float64/vec2"
)

func TestCurve_interface(t *testing.T) {
	var curve T = NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
	if curve == nil {
		t.Errorf("curve does not implement interface T")
	}
}

func TestCurveBBox(t *testing.T) {
	v := NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
	want := vec2.Rect{Min: vec2.T{0, 0}, Max: vec2.T{2, 0.75}}
	got := v.BBox()
	if got.Min != want.Min {
		t.Errorf("BBox Min failed: got %v, want %v", got.Min, want.Min)
	}
	if got.Max != want.Max {
		t.Errorf("BBox Max failed: got %v, want %v", got.Max, want.Max)
	}
}

func TestCurveAt(t *testing.T) {
	v := NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
	want := vec2.T{1, 0.75}
	got := v.At(0.5)
	if got != want {
		t.Errorf("At failed: got %v, want %v", got, want)
	}
}

func TestCurveTangent(t *testing.T) {
	v := NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
	want := vec2.T{3, 0}
	got := v.Tangent(0.5)
	if got != want {
		t.Errorf("Tangent failed: got %v, want %v", got, want)
	}
}

func TestCurveNTangent(t *testing.T) {
	v := NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
	want := vec2.T{1, 0}
	got := v.NTangent(0.5)
	e := 1e-15
	if math.Abs(got[0]-want[0]) > e || math.Abs(got[1]-want[1]) > e {
		t.Errorf("NTangent failed: got %v, want %v", got, want)
	}
}

func TestCurveNormal(t *testing.T) {
	v := NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
	want := vec2.T{0, 3}
	got := v.Normal(0.5)
	if got != want {
		t.Errorf("Normal failed: got %v, want %v", got, want)
	}
}

func TestCurveNNormal(t *testing.T) {
	v := NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
	want := vec2.T{0, 1}
	got := v.NNormal(0.5)
	e := 1e-15
	if math.Abs(got[0]-want[0]) > e || math.Abs(got[1]-want[1]) > e {
		t.Errorf("NNormal failed: got %v, want %v", got, want)
	}
}

func TestCurveSubdivide(t *testing.T) {
	v := NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
	want := []float64{0, 0.125, 0.25, 0.5, 0.75, 0.875, 1}
	got := v.Subdivide(45)
	if len(got) != len(want) {
		t.Errorf("Subdivide failed: got %v, want %v", got, want)
	}
}

// func TestCurveWall(t *testing.T) {
// 	v := NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
// 	want := []Triangle3D{
// 		{vec3.T{0, 0, 0}, vec3.T{1, 0.75, 4}, vec3.T{0, 0, 4}},
// 		{vec3.T{0, 0, 0}, vec3.T{1, 0.75, 0}, vec3.T{1, 0.75, 4}},
// 	}
// 	got := v.Wall(4, 1)
// 	for i, tri := range g {
// 		if tri[0] != want[i][0] || tri[1] != want[i][1] || tri[2] != want[i][2] {
// 			t.Errorf("Wall #%v failed: got %v, want %v", i, tri, want[i])
// 		}
// 	}
// }
//
// func TestCurveBevel(t *testing.T) {
// 	v := NewCurve(vec2.T{0, 0}, vec2.T{0, 1}, vec2.T{2, 1}, vec2.T{2, 0})
// 	want := []Triangle3D{
// 		{vec3.T{0, 0, 4}, vec3.T{1, 1, 5}, vec3.T{0, 1, 5}},
// 		{vec3.T{0, 0, 4}, vec3.T{1, 0, 4}, vec3.T{1, 1, 5}},
// 	}
// 	got := v.Bevel(4, 1, 45, 1)
// 	for i, tri := range g {
// 		if tri[0] != want[i][0] || tri[1] != want[i][1] || tri[2] != want[i][2] {
// 			t.Errorf("NNormal #%v failed: got %v, want %v", i, tri, want[i])
// 		}
// 	}
// }
