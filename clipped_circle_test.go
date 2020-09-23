package geojson

import (
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestClippedCircleNew(t *testing.T) {
	circle := NewCircle(P(-112, 33), 123456.654321, 64)
	clipper := RO(-113, 32.5, -112, 33.5)
	g := NewClippedCircle(circle, clipper, nil)
	exectedJson := `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[-112,33]},"properties":{"type":"Circle","radius":123456.654321,"radius_units":"m"}}{"type":"Polygon","coordinates":[[[-113,32.5],[-112,32.5],[-112,33.5],[-113,33.5],[-113,32.5]]]}}`
	expect(t, g.JSON() == exectedJson)
}


func TestPointClippedCircle(t *testing.T) {
	circle := NewCircle(P(-112, 33), 123456.654321, 64)
	clipper := RO(-113, 32.5, -112, 33.5)
	g := NewClippedCircle(circle, clipper, nil)
	p := NewPoint(geometry.Point{X: -112.26, Y: 33.49})
	expect(t, p.Within(g))
	expect(t, g.Contains(p))
}

func TestClippedCircleContains(t *testing.T) {
	circle := NewCircle(P(-122.4412, 37.7335), 1000, 64)
	clipper := RO(-130, 32.5, -120, 43.5)
	g := NewClippedCircle(circle, clipper, nil)
	expect(t, g.Contains(PO(-122.4412, 37.7335)))
	expect(t, g.Contains(PO(-122.44121, 37.7335)))
	expect(t, g.Contains(
		MPO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.4408378, 37.733)})))
	expect(t, g.Contains(
		NewCircle(P(-122.44121, 37.7335), 500, 64)))
	expect(t, g.Contains(
		LO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.4408378, 37.733)})))
	expect(t, g.Contains(
		MLO([]*geometry.Line{
			L([]geometry.Point{
				P(-122.4408378, 37.7341129),
				P(-122.4408378, 37.733),
			}),
			L([]geometry.Point{
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
			})})))
	expect(t, g.Contains(
		PPO(
			[]geometry.Point{
				P(-122.4408378, 37.7341129),
				P(-122.4408378, 37.733),
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
				P(-122.4408378, 37.7341129),
			},
			[][]geometry.Point{})))

	// Does-not-contain
	expect(t, !g.Contains(PO(-122.265, 37.826)))
	expect(t, !g.Contains(
		NewCircle(P(-122.265, 37.826), 100, 64)))
	expect(t, !g.Contains(
		LO([]geometry.Point{
			P(-122.265, 37.826),
			P(-122.210, 37.860)})))
	expect(t, !g.Contains(
		MPO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.198181, 37.7490)})))
	expect(t, !g.Contains(
		MLO([]*geometry.Line{
			L([]geometry.Point{
				P(-122.265, 37.826),
				P(-122.265, 37.860),
			}),
			L([]geometry.Point{
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
			})})))
	expect(t, !g.Contains(PPO(
		[]geometry.Point{
			P(-122.265, 37.826),
			P(-122.265, 37.860),
			P(-122.210, 37.860),
			P(-122.210, 37.826),
			P(-122.265, 37.826),
		},
		[][]geometry.Point{})))
}

func TestClippedCircleIntersects(t *testing.T) {
	circle := NewCircle(P(-122.4412, 37.7335), 1000, 64)
	clipper := RO(-130, 32.5, -120, 43.5)
	g := NewClippedCircle(circle, clipper, nil)
	expect(t, g.Intersects(PO(-122.4412, 37.7335)))
	expect(t, g.Intersects(PO(-122.44121, 37.7335)))
	expect(t, g.Intersects(
		MPO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.4408378, 37.733)})))
	expect(t, g.Intersects(
		NewCircle(P(-122.44121, 37.7335), 500, 64)))
	expect(t, g.Intersects(
		LO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.4408378, 37.733)})))
	expect(t, g.Intersects(
		LO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.265, 37.826)})))
	expect(t, g.Intersects(
		MLO([]*geometry.Line{
			L([]geometry.Point{
				P(-122.4408378, 37.7341129),
				P(-122.265, 37.826),
			}),
			L([]geometry.Point{
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
			})})))
	expect(t, g.Intersects(
		PPO(
			[]geometry.Point{
				P(-122.4408378, 37.7341129),
				P(-122.265, 37.860),
				P(-122.210, 37.826),
				P(-122.44, 37.7341129),
				P(-122.4408378, 37.7341129),
			},
			[][]geometry.Point{})))
	expect(t, g.Intersects(
		MPO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.198181, 37.7490)})))
	expect(t, g.Intersects(
		MLO([]*geometry.Line{
			L([]geometry.Point{
				P(-122.265, 37.826),
				P(-122.265, 37.860),
			}),
			L([]geometry.Point{
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
			})})))

	// Does-not-intersect
	expect(t, !g.Intersects(PO(-122.265, 37.826)))
	expect(t, !g.Intersects(
		NewCircle(P(-122.265, 37.826), 100, 64)))
	expect(t, !g.Intersects(
		LO([]geometry.Point{
			P(-122.265, 37.826),
			P(-122.210, 37.860)})))
}
