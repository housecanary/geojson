package geom

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tidwall/lotsa"
)

func TestRingScan(t *testing.T) {
	test := func(t *testing.T, index bool) {
		rectangleRing := newRing(rectangle)
		if !index {
			rectangleRing.(*Series).tree = nil
		} else {
			rectangleRing.(*Series).buildTree()
		}
		var segs []Segment
		rectangleRing.ForEachSegment(func(seg Segment, idx int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect := []Segment{
			S(0, 0, 10, 0),
			S(10, 0, 10, 10),
			S(10, 10, 0, 10),
			S(0, 10, 0, 0),
		}
		expect(t, len(segs) == len(segsExpect))
		for i := 0; i < len(segs); i++ {
			expect(t, segs[i] == segsExpect[i])
		}

		segs = nil
		notClosedRing := newRing(rectangle)
		if !index {
			notClosedRing.(*Series).tree = nil
		} else {
			notClosedRing.(*Series).buildTree()
		}
		notClosedRing.ForEachSegment(func(seg Segment, idx int) bool {
			segs = append(segs, seg)
			return true
		})
		expect(t, len(segs) == len(segsExpect))
		for i := 0; i < len(segs); i++ {
			expect(t, segs[i] == segsExpect[i])
		}
	}
	t.Run("Indexed", func(t *testing.T) {
		test(t, true)
	})
	t.Run("Simple", func(t *testing.T) {
		test(t, false)
	})
}

func TestRingSearch(t *testing.T) {
	test := func(t *testing.T, index bool) {
		octagonRing := newRing(octagon)
		if !index {
			octagonRing.(*Series).tree = nil
		} else {
			octagonRing.(*Series).buildTree()
		}
		var segs []Segment
		octagonRing.Search(R(0, 0, 0, 0), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect := []Segment{
			S(0, 3, 3, 0),
		}
		expect(t, checkSegsDups(segsExpect, segs))
		segs = nil
		octagonRing.Search(R(0, 0, 0, 10), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect = []Segment{
			S(3, 10, 0, 7),
			S(0, 7, 0, 3),
			S(0, 3, 3, 0),
		}
		expect(t, checkSegsDups(segsExpect, segs))
		segs = nil
		octagonRing.Search(R(0, 0, 5, 10), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect = []Segment{
			S(3, 0, 7, 0),
			S(7, 10, 3, 10),
			S(3, 10, 0, 7),
			S(0, 7, 0, 3),
			S(0, 3, 3, 0),
		}
		expect(t, checkSegsDups(segsExpect, segs))
	}
	t.Run("Indexed", func(t *testing.T) {
		test(t, true)
	})
	t.Run("Simple", func(t *testing.T) {
		test(t, false)
	})
}

func TestRingIntersectsSegment(t *testing.T) {
	simple := newRing(concave1)
	simple.(*Series).tree = nil
	tree := newRing(concave1)
	tree.(*Series).buildTree()

	expect(t, !ringIntersectsSegment(simple, S(0, 0, 3, 3), true))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 3, 3), true))
	expect(t, !ringIntersectsSegment(simple, S(0, 0, 3, 3), false))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 3, 3), false))

	expect(t, ringIntersectsSegment(simple, S(0, 0, 5, 5), true))
	expect(t, ringIntersectsSegment(tree, S(0, 0, 5, 5), true))
	expect(t, !ringIntersectsSegment(simple, S(0, 0, 5, 5), false))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 5, 5), false))

	expect(t, ringIntersectsSegment(simple, S(0, 0, 10, 10), true))
	expect(t, ringIntersectsSegment(tree, S(0, 0, 10, 10), true))
	expect(t, !ringIntersectsSegment(simple, S(0, 0, 10, 10), false))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 10, 10), false))

	expect(t, ringIntersectsSegment(simple, S(0, 0, 11, 11), true))
	expect(t, ringIntersectsSegment(tree, S(0, 0, 11, 11), true))
	expect(t, !ringIntersectsSegment(simple, S(0, 0, 11, 11), false))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 11, 11), false))

}

func TestRingIntersectsRing(t *testing.T) {
	simple := newRing(concave1)
	simple.(*Series).tree = nil
	tree := newRing(concave1)
	tree.(*Series).buildTree()
	small := newRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small.(*Series).tree = nil

	intersects := func(ring Ring) bool {
		tt := ringIntersectsRing(simple, ring, true)
		if ringIntersectsRing(tree, ring, true) != tt {
			panic("structure mismatch")
		}
		return tt
	}

	intersectsOnEdgeNotAllowed := func(ring Ring) bool {
		tt := ringIntersectsRing(simple, ring, false)
		if ringIntersectsRing(tree, ring, false) != tt {
			panic("structure mismatch")
		}
		return tt
	}

	expect(t, intersects(small))
	expect(t, intersects(small.(*Series).move(-6, 0)))
	expect(t, intersects(small.(*Series).move(6, 0)))
	expect(t, !intersects(small.(*Series).move(-7, 0)))
	expect(t, !intersects(small.(*Series).move(7, 0)))
	expect(t, intersects(small.(*Series).move(1, 1)))
	expect(t, intersects(small.(*Series).move(-1, -1)))
	expect(t, intersects(small.(*Series).move(2, 2)))
	expect(t, !intersects(small.(*Series).move(-2, -2)))
	expect(t, intersects(small.(*Series).move(0, -6)))
	expect(t, intersects(small.(*Series).move(0, 6)))
	expect(t, !intersects(small.(*Series).move(0, -7)))
	expect(t, !intersects(small.(*Series).move(0, 7)))

	expect(t, intersectsOnEdgeNotAllowed(small.(*Series).move(-5, 0)))
	expect(t, intersectsOnEdgeNotAllowed(small.(*Series).move(5, 0)))
	expect(t, intersectsOnEdgeNotAllowed(small.(*Series).move(0, -5)))
	expect(t, intersectsOnEdgeNotAllowed(small.(*Series).move(0, 5)))

	expect(t, !intersectsOnEdgeNotAllowed(small.(*Series).move(-6, 0)))
	expect(t, !intersectsOnEdgeNotAllowed(small.(*Series).move(6, 0)))
	expect(t, !intersectsOnEdgeNotAllowed(small.(*Series).move(0, -6)))
	expect(t, !intersectsOnEdgeNotAllowed(small.(*Series).move(0, 6)))

	expect(t, intersectsOnEdgeNotAllowed(small.(*Series).move(1, 1)))
	expect(t, !intersectsOnEdgeNotAllowed(small.(*Series).move(-1, -1)))

}

func TestBigRandomPIP(t *testing.T) {
	simple := newRing(az)
	simple.(*Series).tree = nil
	tree := newRing(az)
	tree.(*Series).buildTree()
	expect(t, simple.Rect() == tree.Rect())
	rect := tree.Rect()
	start := time.Now()
	for time.Since(start) < time.Second/4 {
		point := P(
			rand.Float64()*(rect.Max.X-rect.Min.X)+rect.Min.X,
			rand.Float64()*(rect.Max.Y-rect.Min.Y)+rect.Min.Y,
		)
		expect(t, ringContainsPoint(tree, point, true) ==
			ringContainsPoint(simple, point, true))
	}
}

func testBig(
	t *testing.T, label string, points []Point, pointIn, pointOut Point,
) {
	N := 10000
	simple := newRing(points)
	simple.(*Series).tree = nil
	tree := newRing(points)
	tree.(*Series).buildTree()
	pointOn := points[len(points)/2]

	expect(t, ringContainsPoint(simple, pointIn, true))
	expect(t, ringContainsPoint(tree, pointIn, true))

	expect(t, ringContainsPoint(simple, pointOn, true))
	expect(t, ringContainsPoint(tree, pointOn, true))

	expect(t, !ringContainsPoint(simple, pointOn, false))
	expect(t, !ringContainsPoint(tree, pointOn, false))

	expect(t, !ringContainsPoint(simple, pointOut, true))
	expect(t, !ringContainsPoint(tree, pointOut, true))
	if os.Getenv("PIPBENCH") == "1" {
		lotsa.Output = os.Stderr
		fmt.Printf(label + "/simp/in  ")
		lotsa.Ops(N, 1, func(_, _ int) {
			ringContainsPoint(simple, pointIn, true)
		})
		fmt.Printf(label + "/tree/in  ")
		lotsa.Ops(N, 1, func(_, _ int) {
			ringContainsPoint(tree, pointIn, true)
		})
		fmt.Printf(label + "/simp/on  ")
		lotsa.Ops(N, 1, func(_, _ int) {
			ringContainsPoint(simple, pointOn, true)
		})
		fmt.Printf(label + "/tree/on  ")
		lotsa.Ops(N, 1, func(_, _ int) {
			ringContainsPoint(tree, pointOn, true)
		})
		fmt.Printf(label + "/simp/out ")
		lotsa.Ops(N, 1, func(_, _ int) {
			ringContainsPoint(simple, pointOut, true)
		})
		fmt.Printf(label + "/tree/out ")
		lotsa.Ops(N, 1, func(_, _ int) {
			ringContainsPoint(tree, pointOut, true)
		})
	}
}

func TestBigArizona(t *testing.T) {
	testBig(t, "az", az, P(-112, 33), P(-114.477539062, 33.99802726))
}

func TestBigTexas(t *testing.T) {
	testBig(t, "tx", tx, P(-98.52539, 29.363027), P(-101.953125, 29.324720161))
}

func TestBigCircle(t *testing.T) {
	circle := CircleRing(P(-100.1, 31.2), 660000, 10000).(*Series).Points()
	if false {
		s := `{"type":"Polygon","coordinates":[[`
		for i, p := range circle {
			if i > 0 {
				s += ","
			}
			s += fmt.Sprintf("[%v,%v]", p.X, p.Y)
		}
		s += `]]}`
		println(s)
	}
	testBig(t, "circ", circle, P(-98.52, 29.363), P(-107.8857, 31.5410))
	circle = CircleRing(P(-100.1, 31.2), 660000, 2).(*Series).Points()
	expect(t, len(circle) == 4)
}

func TestRingContainsRing(t *testing.T) {
	simple := newRing(concave1)
	simple.(*Series).tree = nil
	tree := newRing(concave1)
	tree.(*Series).buildTree()

	expect(t, ringContainsRing(simple, simple, true))
	expect(t, ringContainsRing(simple, tree, true))
	expect(t, ringContainsRing(tree, simple, true))
	expect(t, ringContainsRing(tree, tree, true))

	expect(t, !ringContainsRing(simple, simple, false))
	expect(t, !ringContainsRing(simple, tree, false))
	expect(t, !ringContainsRing(tree, simple, false))
	expect(t, !ringContainsRing(tree, tree, false))

	small := newRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small.(*Series).tree = nil

	expect(t, !ringContainsRing(simple, small, true))
	expect(t, !ringContainsRing(tree, small, true))

	for x := 1.0; x <= 4; x++ {
		expect(t, ringContainsRing(simple, small.(*Series).move(x, 0), true))
		expect(t, ringContainsRing(tree, small.(*Series).move(x, 0), true))
	}
	expect(t, !ringContainsRing(simple, small.(*Series).move(4, 0), false))
	expect(t, !ringContainsRing(tree, small.(*Series).move(4, 0), false))
	for y := 1.0; y <= 4; y++ {
		expect(t, ringContainsRing(simple, small.(*Series).move(0, y), true))
		expect(t, ringContainsRing(tree, small.(*Series).move(0, y), true))
	}
	expect(t, !ringContainsRing(simple, small.(*Series).move(0, 4), false))
	expect(t, !ringContainsRing(tree, small.(*Series).move(0, 4), false))

	for x := -1.0; x >= -4; x-- {
		expect(t, !ringContainsRing(simple, small.(*Series).move(x, 0), true))
		expect(t, !ringContainsRing(tree, small.(*Series).move(x, 0), true))
	}
	expect(t, !ringContainsRing(simple, small.(*Series).move(-4, 0), false))
	expect(t, !ringContainsRing(tree, small.(*Series).move(-4, 0), false))
	for y := -1.0; y >= -4; y-- {
		expect(t, !ringContainsRing(simple, small.(*Series).move(0, y), true))
		expect(t, !ringContainsRing(tree, small.(*Series).move(0, y), true))
	}
	expect(t, !ringContainsRing(simple, small.(*Series).move(0, -4), false))
	expect(t, !ringContainsRing(tree, small.(*Series).move(0, -4), false))

	expect(t, !ringContainsRing(simple, small.(*Series).move(1, 0), false))
	expect(t, !ringContainsRing(tree, small.(*Series).move(1, 0), false))
	expect(t, ringContainsRing(simple, small.(*Series).move(2, 0), false))
	expect(t, ringContainsRing(tree, small.(*Series).move(2, 0), false))
	expect(t, ringContainsRing(simple, small.(*Series).move(2, 2), false))
	expect(t, ringContainsRing(tree, small.(*Series).move(2, 2), false))
	expect(t, !ringContainsRing(simple, small.(*Series).move(-2, -2), false))
	expect(t, !ringContainsRing(tree, small.(*Series).move(-2, -2), false))

	expect(t, !ringContainsRing(simple, small.(*Series).move(5, 0), true))
	expect(t, !ringContainsRing(tree, small.(*Series).move(5, 0), true))
	expect(t, !ringContainsRing(simple, small.(*Series).move(-5, 0), true))
	expect(t, !ringContainsRing(tree, small.(*Series).move(-5, 0), true))

	expect(t, !ringContainsRing(simple, small.(*Series).move(0, 5), true))
	expect(t, !ringContainsRing(tree, small.(*Series).move(0, 5), true))
	expect(t, !ringContainsRing(simple, small.(*Series).move(0, -5), true))
	expect(t, !ringContainsRing(tree, small.(*Series).move(0, -5), true))

}
func TestBowtie(t *testing.T) {
	simple := newRing(bowtie)
	simple.(*Series).tree = nil
	tree := newRing(bowtie)
	tree.(*Series).buildTree()
	square := newRing([]Point{P(3, 3), P(7, 3), P(7, 7), P(3, 7), P(3, 3)})
	square.(*Series).tree = nil

	expect(t, ringIntersectsRing(simple, square, true))
	expect(t, ringIntersectsRing(tree, square, true))
	expect(t, !ringContainsRing(simple, square, true))
	expect(t, !ringContainsRing(tree, square, true))

}

func TestRingVarious(t *testing.T) {
	ring := newRing(octagon[:len(octagon)-1])
	ring.(*Series).buildTree()
	n := 0
	ring.Search(R(0, 0, 10, 10), func(seg Segment, index int) bool {
		n++
		return true
	})
	expect(t, n == 8)
	n = 0
	ring.ForEachSegment(func(seg Segment, idx int) bool {
		n++
		return true
	})
	expect(t, n == 8)
	n = 0
	ring.ForEachSegment(func(seg Segment, idx int) bool {
		n++
		return false
	})
	expect(t, n == 1)
	expect(t, ringIntersectsSegment(ring, S(0, 0, 4, 4), true))
	expect(t, !newRingSimple2([]Point{}).Convex())
	expect(t, newRingSimple2(octagon).Convex())
	expect(t, !newRingIndexed2([]Point{}).Convex())
	expect(t, newRingIndexed2(octagon).Convex())

	ring = newRing(octagon[:len(octagon)-1])
	ring.(*Series).tree = nil
	n = 0
	ring.Search(R(0, 0, 10, 10), func(seg Segment, index int) bool {
		n++
		return false
	})
	expect(t, n == 1)
	n = 0
	ring.ForEachSegment(func(seg Segment, idx int) bool {
		n++
		return true
	})
	expect(t, n == 8)
	expect(t, ringIntersectsSegment(ring, S(0, 0, 4, 4), true))

	small := newRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small.(*Series).tree = nil

	expect(t, ringIntersectsRing(small, ring, true))
	expect(t, ringIntersectsRing(ring, small, true))

	expect(t, raycast(P(0, 0), P(0, 0), P(0, 0)).on)

	ring1 := newRing(octagon)
	ring1.(*Series).tree = nil
	n1 := 0
	ring1.ForEachSegment(func(seg Segment, idx int) bool {
		n1++
		return true
	})
	expect(t, ring1.(*Series).Closed())
	ring2 := newRing(octagon[:len(octagon)-1])
	ring2.(*Series).tree = nil
	n2 := 0
	ring2.ForEachSegment(func(seg Segment, idx int) bool {
		n2++
		return true
	})
	expect(t, n1 == n2)
	expect(t, ring2.(*Series).Closed())

	ring1 = newRing(octagon)
	ring1.(*Series).buildTree()
	n1 = 0
	ring1.ForEachSegment(func(seg Segment, idx int) bool {
		n1++
		return true
	})
	expect(t, ring1.(*Series).Closed())
	ring2 = newRing(octagon[:len(octagon)-1])
	ring2.(*Series).buildTree()
	n2 = 0
	ring2.ForEachSegment(func(seg Segment, idx int) bool {
		n2++
		return true
	})
	expect(t, n1 == n2)
	expect(t, ring2.(*Series).Closed())

}

func newRingSimple2(points []Point) Ring {
	ring := newRing(points)
	ring.(*Series).tree = nil
	return ring
}
func newRingIndexed2(points []Point) Ring {
	ring := newRing(points)
	ring.(*Series).buildTree()
	return ring
}

func TestRingContainsPoint(t *testing.T) {
	expect(t, ringIntersectsPoint(newRingSimple2(octagon), P(4, 4), true))
	expect(t, ringIntersectsPoint(newRingIndexed2(octagon), P(4, 4), true))
}

func TestRingContainsSegment(t *testing.T) {
	expect(t, ringContainsSegment(newRingSimple2(octagon), S(4, 4, 6, 6), true))
	expect(t, ringContainsSegment(newRingIndexed2(octagon), S(4, 4, 6, 6), true))
	expect(t, !ringContainsSegment(newRingSimple2(octagon), S(9, 4, 11, 6), true))
	expect(t, !ringContainsSegment(newRingIndexed2(octagon), S(9, 4, 11, 6), true))
	expect(t, !ringContainsSegment(newRingSimple2(octagon), S(11, 4, 9, 6), true))
	expect(t, !ringContainsSegment(newRingIndexed2(octagon), S(11, 4, 9, 6), true))
	expect(t, !ringContainsSegment(newRingSimple2(concave1), S(11, 4, 9, 6), true))
	expect(t, !ringContainsSegment(newRingIndexed2(concave1), S(11, 4, 9, 6), true))
	expect(t, ringContainsSegment(newRingSimple2(concave1), S(6, 6, 8, 8), true))
	expect(t, ringContainsSegment(newRingIndexed2(concave1), S(6, 6, 8, 8), true))
	expect(t, !ringContainsSegment(newRingSimple2(concave1), S(1, 6, 6, 1), true))
	expect(t, !ringContainsSegment(newRingIndexed2(concave1), S(1, 6, 6, 1), true))
}
func TestRingContainsRect(t *testing.T) {
	expect(t, ringContainsRect(newRingSimple2(octagon), R(4, 4, 6, 6), true))
	expect(t, ringContainsRect(newRingIndexed2(octagon), R(4, 4, 6, 6), true))
	expect(t, ringContainsRect(newRingSimple2(octagon), R(4, 4, 6, 6), false))
	expect(t, ringContainsRect(newRingIndexed2(octagon), R(4, 4, 6, 6), false))
}
func TestRingIntersectsRect(t *testing.T) {
	expect(t, ringIntersectsRect(newRingSimple2(octagon), R(9, 4, 11, 6), true))
	expect(t, ringIntersectsRect(newRingIndexed2(octagon), R(9, 4, 11, 6), true))
	expect(t, !ringIntersectsRect(newRingSimple2(octagon), R(10, 4, 12, 6), false))
	expect(t, !ringIntersectsRect(newRingIndexed2(octagon), R(10, 4, 12, 6), false))
	expect(t, ringIntersectsRect(newRingSimple2(octagon), R(10, 4, 12, 6), true))
	expect(t, ringIntersectsRect(newRingIndexed2(octagon), R(10, 4, 12, 6), true))
	expect(t, !ringIntersectsRect(newRingSimple2(octagon), R(11, 4, 13, 6), true))
	expect(t, !ringIntersectsRect(newRingIndexed2(octagon), R(11, 4, 13, 6), true))
}
func TestRingContainsPoly(t *testing.T) {
	expect(t, ringContainsPoly(newRingSimple2(octagon), NewPoly(octagon, nil), true))
	expect(t, ringContainsPoly(newRingIndexed2(octagon), NewPoly(octagon, nil), true))
	expect(t, !ringContainsPoly(newRingSimple2(octagon), NewPoly(octagon, nil), false))
	expect(t, !ringContainsPoly(newRingIndexed2(octagon), NewPoly(octagon, nil), false))
}
func TestRingIntersectsPoly(t *testing.T) {
	expect(t, ringIntersectsPoly(newRingSimple2(octagon).(*Series).move(5, 0), NewPoly(octagon, nil), true))
	expect(t, ringIntersectsPoly(newRingIndexed2(octagon).(*Series).move(5, 0), NewPoly(octagon, nil), true))
	expect(t, ringIntersectsPoly(newRingSimple2(octagon).(*Series).move(10, 0), NewPoly(octagon, nil), true))
	expect(t, ringIntersectsPoly(newRingIndexed2(octagon).(*Series).move(10, 0), NewPoly(octagon, nil), true))
	expect(t, !ringIntersectsPoly(newRingSimple2(octagon).(*Series).move(10, 0), NewPoly(octagon, nil), false))
	expect(t, !ringIntersectsPoly(newRingIndexed2(octagon).(*Series).move(10, 0), NewPoly(octagon, nil), false))

	expect(t, !ringIntersectsPoly(newRingIndexed2(
		[]Point{P(4, 4), P(6, 4), P(6, 6), P(4, 6), P(4, 4)},
	), NewPoly(
		[]Point{P(0, 0), P(10, 0), P(10, 10), P(0, 10), P(0, 0)},
		[][]Point{{P(3, 3), P(7, 3), P(7, 7), P(3, 7), P(3, 3)}},
	), false))

}

func TestSegmentsIntersect(t *testing.T) {
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(11, 0), P(21, 10)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(-11, 0), P(-21, 10)))
	expect(t, !segmentsIntersect(P(10, 10), P(0, 10), P(11, 0), P(21, 10)))
	expect(t, !segmentsIntersect(P(10, 10), P(0, 10), P(-11, 0), P(-21, 10)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(0, 11), P(10, 21)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(0, -11), P(10, -21)))
	expect(t, !segmentsIntersect(P(10, 10), P(0, 0), P(0, 11), P(10, 21)))
	expect(t, !segmentsIntersect(P(10, 10), P(0, 0), P(0, -11), P(10, -21)))

	expect(t, !segmentsIntersect(P(0, 0), P(10, 0), P(11, 0), P(21, 0)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 0), P(0, 1), P(10, 1)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 0), P(0, -1), P(10, -1)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(1, 0), P(11, 10)))

}

func BenchmarkCircleRect(b *testing.B) {
	for i := 4; i < 256; i *= 2 {
		indexed := CircleRing(P(-112, 33), 1000, i)
		indexed.(*Series).buildTree()
		simple := CircleRing(P(-112, 33), 1000, i)
		simple.(*Series).tree = nil
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			b.Run("Simple", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					simple.Rect()
				}
			})
			b.Run("Indexed", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					indexed.Rect()
				}
			})
		})
	}
}