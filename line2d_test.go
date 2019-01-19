package parametric2d

import (
	"math"
	"testing"

	"github.com/gmlewis/go3d/float64/vec2"
	"github.com/gmlewis/go3d/float64/vec3"
)

func TestLine_interface(t *testing.T) {
	var line T = NewLine(vec2.T{0, 1}, vec2.T{2, -1})
	if line == nil {
		t.Errorf("line does not implement interface T")
	}
}

func TestBBox(t *testing.T) {
	v := NewLine(vec2.T{0, 1}, vec2.T{2, -1})
	want := vec2.Rect{Min: vec2.T{0, -1}, Max: vec2.T{2, 1}}
	got := v.BBox()
	if got.Min[0] != want.Min[0] || got.Min[1] != want.Min[1] {
		t.Errorf("BBox Min failed: got %v, want %v", got.Min, want.Min)
	}
	if got.Max[0] != want.Max[0] || got.Max[1] != want.Max[1] {
		t.Errorf("BBox Max failed: got %v, want %v", got.Max, want.Max)
	}
}

func TestAt(t *testing.T) {
	v := NewLine(vec2.T{0, 1}, vec2.T{2, -1})
	want := vec2.T{1, 0}
	got := v.At(0.5)
	if got[0] != want[0] || got[1] != want[1] {
		t.Errorf("At failed: got %v, want %v", got, want)
	}
}

func TestTangent(t *testing.T) {
	v := NewLine(vec2.T{0, 1}, vec2.T{2, -1})
	want := vec2.T{2, -2}
	got := v.Tangent(0.5)
	if got[0] != want[0] || got[1] != want[1] {
		t.Errorf("Tangent failed: got %v, want %v", got, want)
	}
}

func TestNTangent(t *testing.T) {
	v := NewLine(vec2.T{0, 1}, vec2.T{2, -1})
	want := vec2.T{0.5 * math.Sqrt(2), -0.5 * math.Sqrt(2)}
	got := v.NTangent(0.5)
	e := 1e-15
	if math.Abs(got[0]-want[0]) > e || math.Abs(got[1]-want[1]) > e {
		t.Errorf("NTangent failed: got %v, want %v", got, want)
	}
}

func TestNormal(t *testing.T) {
	v := NewLine(vec2.T{0, 1}, vec2.T{2, -1})
	want := vec2.T{2, 2}
	got := v.Normal(0.5)
	if got[0] != want[0] || got[1] != want[1] {
		t.Errorf("Normal failed: got %v, want %v", got, want)
	}
}

func TestNNormal(t *testing.T) {
	v := NewLine(vec2.T{0, 1}, vec2.T{2, -1})
	want := vec2.T{0.5 * math.Sqrt(2), 0.5 * math.Sqrt(2)}
	got := v.NNormal(0.5)
	e := 1e-15
	if math.Abs(got[0]-want[0]) > e || math.Abs(got[1]-want[1]) > e {
		t.Errorf("NNormal failed: got %v, want %v", got, want)
	}
}

func TestWall(t *testing.T) {
	v := NewLine(vec2.T{0, 0}, vec2.T{1, 0})
	want := []Triangle3D{
		{vec3.T{0, 0, 0}, vec3.T{1, 0, 4}, vec3.T{0, 0, 4}},
		{vec3.T{0, 0, 0}, vec3.T{1, 0, 0}, vec3.T{1, 0, 4}},
	}
	got, _ := v.Wall(4, 1, false)
	for i, tri := range got {
		if tri[0] != want[i][0] || tri[1] != want[i][1] || tri[2] != want[i][2] {
			t.Errorf("NNormal #%v failed: got %v, want %v", i, tri, want[i])
		}
	}
	got, _ = v.Wall(4, 1, true)
	for i, tri := range got {
		if tri[0] != want[i][0] || tri[1] != want[i][2] || tri[2] != want[i][1] {
			t.Errorf("NNormal #%v failed: got %v, want [%v %v %v]", i, tri, want[i][0], want[i][2], want[i][1])
		}
	}
}

func TestBevel(t *testing.T) {
	v := NewLine(vec2.T{0, 0}, vec2.T{1, 0})
	want := []Triangle3D{
		{vec3.T{0, 0, 4}, vec3.T{1, 1, 5}, vec3.T{0, 1, 5}},
		{vec3.T{0, 0, 4}, vec3.T{1, 0, 4}, vec3.T{1, 1, 5}},
	}
	got, _ := v.Bevel(4, 1, 45, 1, false, &vec2.T{0, 1}, &vec2.T{0, 1})
	for i, tri := range got {
		if tri[0] != want[i][0] || tri[1] != want[i][1] || tri[2] != want[i][2] {
			t.Errorf("NNormal #%v failed: got %v, want %v", i, tri, want[i])
		}
	}
	// w = []Triangle3D{
	// 	{vec3.T{0, 0, 4}, vec3.T{0, -1, 5}, vec3.T{1, -1, 5}},
	// 	{vec3.T{0, 0, 4}, vec3.T{1, -1, 5}, vec3.T{1, 0, 4}},
	// }
	// got = v.Bevel(4, 1, 45, 1, true, &vec2.T{0, 1}, &vec2.T{0, 1})
	// for i, tri := range got {
	// 	if tri[0] != want[i][0] || tri[1] != want[i][1] || tri[2] != want[i][2] {
	// 		t.Errorf("NNormal #%v failed: got %v, want %v", i, tri, want[i])
	// 	}
	// }
}
