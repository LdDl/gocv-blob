package blob

import (
	"fmt"
	"image"
	"math"
	"time"

	uuid "github.com/satori/go.uuid"
	"gocv.io/x/gocv"
)

// SimpleBlobie Simplest blob implementation
type SimpleBlobie struct {
	ID                    uuid.UUID
	CurrentRect           image.Rectangle
	Center                image.Point
	Area                  float64
	Diagonal              float64
	AspectRatio           float64
	Track                 []image.Point
	TrackTime             []time.Time
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

// NewSimpleBlobie - Constructor for SimpleBlobie (default values)
func NewSimpleBlobie(rect image.Rectangle, options *BlobOptions) Blobie {
	center := image.Pt((rect.Min.X*2+rect.Dx())/2, (rect.Min.Y*2+rect.Dy())/2)
	width := float64(rect.Dx())
	height := float64(rect.Dy())
	blobie := SimpleBlobie{
		CurrentRect:         rect,
		Center:              center,
		Area:                width * height,
		Diagonal:            math.Sqrt(math.Pow(width, 2) + math.Pow(height, 2)),
		AspectRatio:         width / height,
		Track:               []image.Point{center},
		isExists:            true,
		isStillBeingTracked: true,
		noMatchTimes:        0,
		crossedLine:         false,
	}
	if options != nil {
		blobie.TrackTime = []time.Time{options.Time}
		blobie.maxPointsInTrack = options.MaxPointsInTrack
		blobie.classID = options.ClassID
		blobie.className = options.ClassName
	} else {
		blobie.TrackTime = []time.Time{time.Now()}
		blobie.maxPointsInTrack = 10
		blobie.classID = -1
		blobie.className = "No class"
	}
	return &blobie
}

// NewBlobieDefaults - Constructor for SimpleBlobie (default values)
//
// Default values are:
// maxPointsInTrack = 10
//
func NewBlobieDefaults(rect image.Rectangle) *SimpleBlobie {
	center := image.Pt((rect.Min.X*2+rect.Dx())/2, (rect.Min.Y*2+rect.Dy())/2)
	width := float64(rect.Dx())
	height := float64(rect.Dy())
	return &SimpleBlobie{
		CurrentRect:         rect,
		Center:              center,
		Area:                width * height,
		Diagonal:            math.Sqrt(math.Pow(width, 2) + math.Pow(height, 2)),
		AspectRatio:         width / height,
		Track:               []image.Point{center},
		TrackTime:           []time.Time{time.Now()},
		maxPointsInTrack:    10,
		isExists:            true,
		isStillBeingTracked: true,
		noMatchTimes:        0,

		classID:     -1,
		className:   "No class",
		crossedLine: false,
	}
}

// PredictNextPosition - Predict next N coordinates
func (b *SimpleBlobie) PredictNextPosition(n int) {
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

// Update - Update info about blob
func (b *SimpleBlobie) Update(newb Blobie) error {
	var newbCast *SimpleBlobie
	switch newb.(type) {
	case *SimpleBlobie:
		newbCast = newb.(*SimpleBlobie)
		break
	default:
		return fmt.Errorf("SimpleBlobie.Update() method must accept interface of type *SimpleBlobie")
	}
	b.CurrentRect = newbCast.CurrentRect
	b.Center = newbCast.Center
	b.Area = newbCast.Area
	b.Diagonal = newbCast.Diagonal
	b.AspectRatio = newbCast.AspectRatio
	b.isStillBeingTracked = true
	b.isExists = true
	// Append new point to track
	b.Track = append(b.Track, newbCast.Center)
	b.TrackTime = append(b.TrackTime, newbCast.TrackTime[len(newbCast.TrackTime)-1])
	// Restrict number of points in track (shift to the left)
	if len(b.Track) > b.maxPointsInTrack {
		b.Track = b.Track[1:]
	}
	return nil
}

func (sb *SimpleBlobie) GetID() uuid.UUID {
	return sb.ID
}

func (sb *SimpleBlobie) GetCenter() image.Point {
	return sb.Center
}

func (sb *SimpleBlobie) GetCurrentRect() image.Rectangle {
	return sb.CurrentRect
}

func (sb *SimpleBlobie) GetTrack() []image.Point {
	return sb.Track
}

func (sb *SimpleBlobie) GetTimestamps() []time.Time {
	return sb.TrackTime
}

func (sb *SimpleBlobie) GetDiagonal() float64 {
	return sb.Diagonal
}

func (sb *SimpleBlobie) GetPredictedNextPosition() image.Point {
	return sb.PredictedNextPosition
}

func (sb *SimpleBlobie) NoMatchTimes() int {
	return sb.noMatchTimes
}

func (sb *SimpleBlobie) Exists() bool {
	return sb.isExists
}

func (sb *SimpleBlobie) SetID(id uuid.UUID) {
	sb.ID = id
}

func (sb *SimpleBlobie) SetTracking(isStillBeingTracked bool) {
	sb.isStillBeingTracked = isStillBeingTracked
}

func (sb *SimpleBlobie) IncrementNoMatchTimes() {
	sb.noMatchTimes++
}

func (sb *SimpleBlobie) SetExists(isExists bool) {
	sb.isExists = isExists
}

// GetClassID Returns class identifier [SimpleBlobie]
func (b *SimpleBlobie) GetClassID() int {
	return b.classID
}

// GetClassName Returns class name [SimpleBlobie]
func (b *SimpleBlobie) GetClassName() string {
	return b.className
}

// SetDraw Sets options for drawing [SimpleBlobie]
func (b *SimpleBlobie) SetDraw(drawOptions *DrawOptions) {
	b.drawingOptions = drawOptions
}

// DrawTrack Draws blob's track [SimpleBlobie]
func (b *SimpleBlobie) DrawTrack(mat *gocv.Mat, optionalText ...string) {
	if b.drawingOptions == nil {
		b.drawingOptions = NewDrawOptionsDefault()
	}
	gocv.Rectangle(mat, b.CurrentRect, b.drawingOptions.BBoxColor.Color, b.drawingOptions.BBoxColor.Thickness)
	if b.isStillBeingTracked {
		for i := range b.Track {
			gocv.Circle(mat, b.Track[i], b.drawingOptions.CentroidColor.Radius, b.drawingOptions.CentroidColor.Color, b.drawingOptions.CentroidColor.Thickness)
		}
		shiftTextY := 0
		for i := 0; i < len(optionalText); i++ {
			text := optionalText[i]
			if text != "" {
				textSize := gocv.GetTextSize(text, b.drawingOptions.TextColor.Font, b.drawingOptions.TextColor.Scale, b.drawingOptions.TextColor.Thickness)
				anchor := image.Pt(b.CurrentRect.Min.X, b.CurrentRect.Min.Y-shiftTextY-b.drawingOptions.BBoxColor.Thickness) // substract extra margin = Thickness of BBox
				textRect := image.Rectangle{Min: image.Point{X: anchor.X, Y: anchor.Y - textSize.Y}, Max: image.Point{X: anchor.X + textSize.X, Y: anchor.Y}}
				gocv.Rectangle(mat, textRect, b.drawingOptions.BBoxColor.Color, b.drawingOptions.BBoxColor.Thickness)
				gocv.PutText(mat, text, anchor, b.drawingOptions.TextColor.Font, b.drawingOptions.TextColor.Scale, b.drawingOptions.TextColor.Color, b.drawingOptions.TextColor.Thickness)
				shiftTextY += textSize.Y
			}
		}
	}
}
