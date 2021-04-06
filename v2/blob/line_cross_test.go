package blob

import (
	"image"
	"testing"
)

func TestHorizontalLineTest(t *testing.T) {
	horizontalLine := [][]int{
		[]int{4, 35},
		[]int{73, 35},
	}
	direction := true // true - TO us
	allblobies := NewBlobiesDefaults()

	simpleB_time0 := NewSimpleBlobie(image.Rect(26, 8, 44, 18), nil)
	simpleB_time1 := NewSimpleBlobie(image.Rect(26, 20, 44, 30), nil)
	simpleB_time2 := NewSimpleBlobie(image.Rect(26, 32, 44, 42), nil)

	allblobies.MatchToExisting([]Blobie{simpleB_time0, simpleB_time1, simpleB_time2})

	for _, b := range allblobies.Objects {
		if b.IsCrossedTheLine(horizontalLine[0][1], horizontalLine[0][0], horizontalLine[1][0], direction) {
			t.Logf("[HORIZONTAL] Correct when direction is TO US")
		} else {
			t.Error("[HORIZONTAL] Incorrect when direction is TO US")
		}

		if b.IsCrossedTheLine(horizontalLine[0][1], horizontalLine[0][0], horizontalLine[1][0], !direction) {
			t.Error("[HORIZONTAL] Incorrect when direction is FROM US")
		} else {
			t.Logf("[HORIZONTAL] Correct when direction is FROM US")
		}
	}
}

func TestObliqueLineTest(t *testing.T) {
	shoudlCrossObliqueLine := [][]int{
		[]int{4, 35},
		[]int{71, 31},
	}
	shoudlNOTCrossObliqueLine := [][]int{
		[]int{4, 35},
		[]int{71, 45},
	}

	direction := true // true - TO us
	allblobies := NewBlobiesDefaults()

	simpleB_time0 := NewSimpleBlobie(image.Rect(26, 8, 44, 18), nil)
	simpleB_time1 := NewSimpleBlobie(image.Rect(26, 20, 44, 30), nil)
	simpleB_time2 := NewSimpleBlobie(image.Rect(26, 32, 44, 42), nil)

	allblobies.MatchToExisting([]Blobie{simpleB_time0, simpleB_time1, simpleB_time2})

	for _, b := range allblobies.Objects {
		if b.IsCrossedTheObliqueLine(shoudlCrossObliqueLine[0][0], shoudlCrossObliqueLine[0][1], shoudlCrossObliqueLine[1][0], shoudlCrossObliqueLine[1][1], direction) {
			t.Logf("[OBLIQUE] Correct when direction is TO US")
		} else {
			t.Error("[OBLIQUE ERR] Incorrect when direction is TO US")
		}

		if b.IsCrossedTheObliqueLine(shoudlCrossObliqueLine[0][0], shoudlCrossObliqueLine[0][1], shoudlCrossObliqueLine[1][0], shoudlCrossObliqueLine[1][1], !direction) {
			t.Error("[OBLIQUE ERR] Incorrect when direction is FROM US")
		} else {
			t.Logf("[OBLIQUE] Correct when direction is FROM US")
		}

		if b.IsCrossedTheObliqueLine(shoudlNOTCrossObliqueLine[0][0], shoudlNOTCrossObliqueLine[0][1], shoudlNOTCrossObliqueLine[1][0], shoudlNOTCrossObliqueLine[1][1], direction) {
			t.Error("[OBLIQUE NOT CROSS ERR] Incorrect when direction is TO US")
		} else {
			t.Logf("[OBLIQUE NOT CROSS] Correct when direction is TO US") 
		}

		if b.IsCrossedTheObliqueLine(shoudlNOTCrossObliqueLine[0][0], shoudlNOTCrossObliqueLine[0][1], shoudlNOTCrossObliqueLine[1][0], shoudlNOTCrossObliqueLine[1][1], !direction) {
			t.Error("[OBLIQUE NOT CROSS ERR] Incorrect when direction is FROM US")
		} else {
			t.Logf("[OBLIQUE NOT CROSS] Correct when direction is FROM US") 
		}
	}
}
