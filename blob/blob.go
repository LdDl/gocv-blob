package blob

import (
	"image"
	"math"
)

// Blobie - Main blob structure
type Blobie struct {
	CurrentRect           image.Rectangle
	Center                image.Point
	Area                  float64
	Diagonal              float64
	AspectRatio           float64
	Track                 []image.Point
	maxPointsInTrack      int
	isExists              bool
	isStillBeingTracked   bool
	noMatchTimes          int
	PredictedNextPosition image.Point

	// For array tracker
	crossedLine bool
}

// NewBlobie - Constructor for Blobie (default values)
func NewBlobie(rect image.Rectangle, maxPointsInTrack int) *Blobie {
	center := image.Pt((rect.Min.X*2+rect.Dx())/2, (rect.Min.Y*2+rect.Dy())/2)
	width := float64(rect.Dx())
	height := float64(rect.Dy())
	return &Blobie{
		CurrentRect:         rect,
		Center:              center,
		Area:                width * height,
		Diagonal:            math.Sqrt(math.Pow(width, 2) + math.Pow(height, 2)),
		AspectRatio:         width / height,
		Track:               []image.Point{center},
		maxPointsInTrack:    maxPointsInTrack,
		isExists:            true,
		isStillBeingTracked: true,
		noMatchTimes:        0,

		crossedLine: false,
	}
}

// NewBlobieDefaults - Constructor for Blobie (default values)
//
// Default values are:
// maxPointsInTrack = 10
//
func NewBlobieDefaults(rect image.Rectangle) *Blobie {
	center := image.Pt((rect.Min.X*2+rect.Dx())/2, (rect.Min.Y*2+rect.Dy())/2)
	width := float64(rect.Dx())
	height := float64(rect.Dy())
	return &Blobie{
		CurrentRect:         rect,
		Center:              center,
		Area:                width * height,
		Diagonal:            math.Sqrt(math.Pow(width, 2) + math.Pow(height, 2)),
		AspectRatio:         width / height,
		Track:               []image.Point{center},
		maxPointsInTrack:    10,
		isExists:            true,
		isStillBeingTracked: true,
		noMatchTimes:        0,

		crossedLine: false,
	}
}

// Update - Update info about blob
func (b *Blobie) Update(newb Blobie) {
	b.CurrentRect = newb.CurrentRect
	b.Center = newb.Center
	b.Area = newb.Area
	b.Diagonal = newb.Diagonal
	b.AspectRatio = newb.AspectRatio
	b.isStillBeingTracked = true
	b.isExists = true
	// Append new point to track
	b.Track = append(b.Track, newb.Center)
	// Restrict number of points in track (shift to the left)
	if len(b.Track) > b.maxPointsInTrack {
		// b.Track = b.Track[1:]
	}
}

// GetLastPoint - Return last point from blob's track
func (b *Blobie) GetLastPoint() image.Point {
	return b.Track[len(b.Track)-1]
}

// PredictNextPosition - Predict next N coordinates
func (b *Blobie) PredictNextPosition(n int) {
	account := min(n, len((*b).Track))
	prev := len((*b).Track) - 1
	current := prev - 1
	var deltaX, deltaY, sum int = 0, 0, 0
	for i := 1; i < int(account); i++ {
		deltaX += (((*b).Track)[current].X - ((*b).Track)[prev].X) * i
		deltaY += (((*b).Track)[current].Y - ((*b).Track)[prev].Y) * i
		sum += i
	}
	if sum > 0 {
		deltaX /= sum
		deltaY /= sum
	}
	(*b).PredictedNextPosition.X = (*b).Track[len((*b).Track)-1].X + deltaX
	(*b).PredictedNextPosition.Y = (*b).Track[len((*b).Track)-1].Y + deltaY
}

// // IsCrossedTheLine - Check if blob crossed the HORIZONTAL line
// func (b *Blobie) IsCrossedTheLine(horizontal int, counter *int, direction bool) bool {
// 	if (*b).isStillBeingTracked == true && len((*b).Track) >= 2 && (*b).counted == false {
// 		prevFrame := len((*b).Track) - 2
// 		currFrame := len((*b).Track) - 1
// 		if direction {
// 			if (*b).Track[prevFrame].Y <= horizontal && (*b).Track[currFrame].Y > horizontal { // TO us
// 				*counter++
// 				b.AsCounted()
// 				return true
// 			}
// 		} else {
// 			if (*b).Track[prevFrame].Y > horizontal && (*b).Track[currFrame].Y <= horizontal { // FROM us
// 				*counter++
// 				b.AsCounted()
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// // DrawTrack - Draw blob's track
// func (b *Blobie) DrawTrack(mat *gocv.Mat, optionalText string) {
// 	if (*b).isStillBeingTracked == true {
// 		for i := range (*b).Track {
// 			gocv.Circle(mat, (*b).Track[i], 4, color.RGBA{255, 0, 0, 0}, 1)
// 		}
// 		gocv.Rectangle(mat, (*b).CurrentRect, color.RGBA{0, 255, 255, 0}, 2)
// 		pt := image.Pt((*b).CurrentRect.Min.X, (*b).CurrentRect.Min.Y)
// 		gocv.PutText(mat, optionalText, pt, gocv.FontHersheyPlain, 1.2, color.RGBA{0, 255, 0, 0}, 2)
// 	}
// }

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// func minf(x, y float32) float32 {
// 	if x < y {
// 		return x
// 	}
// 	return y
// }

// func maxf(x, y float32) float32 {
// 	if x > y {
// 		return x
// 	}
// 	return y
// }

func distanceBetweenBlobies(b1 *Blobie, b2 *Blobie) float64 {
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
