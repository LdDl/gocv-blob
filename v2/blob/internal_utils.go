package blob

import (
	"image"
	"math"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func minf64(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func maxf64(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func distanceBetweenBlobies(b1 *SimpleBlobie, b2 *SimpleBlobie) float64 {
	return distanceBetweenPointsPtr(&b1.Center, &b2.Center)
}

func distanceBetweenPointsPtr(p1 *image.Point, p2 *image.Point) float64 {
	intX := math.Abs(float64(p1.X - p2.X))
	intY := math.Abs(float64(p1.Y - p2.Y))
	return math.Sqrt(math.Pow(intX, 2) + math.Pow(intY, 2))
}

func distanceBetweenPoints(p1 image.Point, p2 image.Point) float64 {
	intX := math.Abs(float64(p1.X - p2.X))
	intY := math.Abs(float64(p1.Y - p2.Y))
	return math.Sqrt(math.Pow(intX, 2) + math.Pow(intY, 2))
}
