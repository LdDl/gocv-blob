package blob

import (
	"image"
	"math"

	uuid "github.com/satori/go.uuid"
	"gocv.io/x/gocv"
)

// Blobie - Main blob structure
type Blobie struct {
	ID                    uuid.UUID
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

	classID   int
	className string

	// For array tracker
	drawingOptions *DrawOptions
	crossedLine    bool
}

// NewBlobie - Constructor for Blobie (default values)
func NewBlobie(rect image.Rectangle, maxPointsInTrack, classID int, className string) *Blobie {
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

		classID:     classID,
		className:   className,
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

		classID:     -1,
		className:   "No class",
		crossedLine: false,
	}
}

// SetClass - Set class information (eg. classID=1, className=vehicle)
func (b *Blobie) SetClass(classID int, className string) {
	b.SetClassID(classID)
	b.SetClassName(className)
}

// SetClassID - Set class identifier
func (b *Blobie) SetClassID(classID int) {
	b.classID = classID
}

// SetClassName - Set class name
func (b *Blobie) SetClassName(className string) {
	b.className = className
}

// GetClassID - Return class identifier
func (b *Blobie) GetClassID() int {
	return b.classID
}

// GetClassName - Return class name
func (b *Blobie) GetClassName() string {
	return b.className
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
		b.Track = b.Track[1:]
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

// DrawTrack - Draw blob's track
func (b *Blobie) DrawTrack(mat *gocv.Mat, optionalText string) {
	if b.drawingOptions == nil {
		b.drawingOptions = NewDrawOptionsDefault()
	}
	gocv.Rectangle(mat, b.CurrentRect, b.drawingOptions.BBoxColor.Color, b.drawingOptions.BBoxColor.Thickness)
	if b.isStillBeingTracked {
		for i := range b.Track {
			gocv.Circle(mat, b.Track[i], b.drawingOptions.CentroidColor.Radius, b.drawingOptions.CentroidColor.Color, b.drawingOptions.CentroidColor.Thickness)
		}
		pt := image.Pt(b.CurrentRect.Min.X, b.CurrentRect.Min.Y)
		gocv.PutText(mat, optionalText, pt, gocv.FontHersheyPlain, b.drawingOptions.TextColor.Scale, b.drawingOptions.TextColor.Color, b.drawingOptions.TextColor.Thickness)
	}
}

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
