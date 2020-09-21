package geojson

import (
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestClipLineStringSimple(t *testing.T) {
	ls := LO([]geometry.Point{
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 1}})
	clipped := Clip(ls, RO(1.5, 0.5, 2.5, 1.8), nil)
	cl, ok := clipped.(*MultiLineString)
	if !ok {
		t.Fatal("wrong type")
	}
	if len(cl.Children()) != 2 {
		t.Fatal("result must have two parts in MultiString")
	}
}

func TestClipPolygonSimple(t *testing.T) {
	exterior := []geometry.Point{
		{X: 2, Y: 2},
		{X: 1, Y: 2},
		{X: 1.5, Y: 1.5},
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 2, Y: 2},
	}
	holes := [][]geometry.Point{
		[]geometry.Point{
			{X: 1.9, Y: 1.9},
			{X: 1.2, Y: 1.9},
			{X: 1.45, Y: 1.65},
			{X: 1.9, Y: 1.5},
			{X: 1.9, Y: 1.9},
		},
	}
	polygon := PPO(exterior, holes)
	clipped := Clip(polygon, RO(1.3, 1.3, 1.4, 2.15), nil)
	cp, ok := clipped.(*Polygon)
	if !ok {
		t.Fatal("wrong type")
	}
	if cp.Base().Exterior.Empty() {
		t.Fatal("Empty result.")
	}
	if len(cp.Base().Holes) != 1 {
		t.Fatal("result must be a two-ring Polygon")
	}
}

func TestClipPolygon2(t *testing.T) {
	exterior := []geometry.Point{
		{X: 2, Y: 2},
		{X: 1, Y: 2},
		{X: 1.5, Y: 1.5},
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 2, Y: 2},
	}
	holes := [][]geometry.Point{
		[]geometry.Point{
			{X: 1.9, Y: 1.9},
			{X: 1.2, Y: 1.9},
			{X: 1.45, Y: 1.65},
			{X: 1.9, Y: 1.5},
			{X: 1.9, Y: 1.9},
		},
	}
	polygon := PPO(exterior, holes)
	clipped := Clip(polygon, RO(1.1, 0.8, 1.15, 2.1), nil)
	cp, ok := clipped.(*Polygon)
	if !ok {
		t.Fatal("wrong type")
	}
	if cp.Base().Exterior.Empty() {
		t.Fatal("Empty result.")
	}
	if len(cp.Base().Holes) != 0 {
		t.Fatal("result must be a single-ring Polygon")
	}
}
