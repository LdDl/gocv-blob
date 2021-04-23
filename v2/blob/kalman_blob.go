package blob

import (
	"fmt"
	"image"
	"math"
	"time"

	kf "github.com/LdDl/kalman-filter"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"

	uuid "github.com/satori/go.uuid"
	"gocv.io/x/gocv"
)

// KalmanBlobie Blob implementation based on Kalman filter.
// For more ref. see: https://en.wikipedia.org/wiki/Kalman_filter
type KalmanBlobie struct {
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

	// Kalman filter wrapping
	pointTracker *kf.PointTracker
	yMatrix      *mat.Dense
	uMatrix      *mat.Dense
	dt           float64

	// For array tracker
	drawingOptions *DrawOptions
	crossedLine    bool
}

// NewKalmanBlobie - Constructor for KalmanBlobie (default values)
func NewKalmanBlobie(rect image.Rectangle, options *BlobOptions) Blobie {
	center := image.Pt((rect.Min.X*2+rect.Dx())/2, (rect.Min.Y*2+rect.Dy())/2)
	width := float64(rect.Dx())
	height := float64(rect.Dy())
	centerX, centerY := float64(center.X), float64(center.Y)
	kalmanBlobie := KalmanBlobie{
		CurrentRect:         rect,
		Center:              center,
		Area:                width * height,
		Diagonal:            math.Sqrt(math.Pow(width, 2) + math.Pow(height, 2)),
		AspectRatio:         width / height,
		Track:               []image.Point{center},
		isExists:            true,
		isStillBeingTracked: true,
		noMatchTimes:        0,
		pointTracker:        kf.NewPointTracker(),
		yMatrix:             mat.NewDense(2, 1, []float64{centerX, centerY}),
		uMatrix:             mat.NewDense(4, 1, []float64{0.0, 0.0, 0.0, 0.0}),
		crossedLine:         false,
	}
	kalmanBlobie.pointTracker.SetStateValue(centerX, centerY, 0, 0)
	if options != nil {
		kalmanBlobie.TrackTime = []time.Time{options.Time}
		kalmanBlobie.maxPointsInTrack = options.MaxPointsInTrack
		kalmanBlobie.classID = options.ClassID
		kalmanBlobie.className = options.ClassName
		kalmanBlobie.dt = options.TimeDeltaSeconds
		kalmanBlobie.pointTracker.SetTime(options.TimeDeltaSeconds)
	} else {
		kalmanBlobie.TrackTime = []time.Time{time.Now()}
		kalmanBlobie.maxPointsInTrack = 10
		kalmanBlobie.classID = -1
		kalmanBlobie.className = "No class"
		kalmanBlobie.dt = 1.0
		kalmanBlobie.pointTracker.SetTime(1.0)
	}
	return &kalmanBlobie
}

// PredictNextPosition - Predict next N coordinates
func (b *KalmanBlobie) PredictNextPosition(n int) {
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
func (b *KalmanBlobie) Update(newb Blobie) error {
	var newbCast *KalmanBlobie
	switch newb.(type) {
	case *KalmanBlobie:
		newbCast = newb.(*KalmanBlobie)
		break
	default:
		return fmt.Errorf("KalmanBlobie.Update() method must accept interface of type *KalmanBlobie")
	}
	newCenterX, newCenterY := float64(newbCast.Center.X), float64(newbCast.Center.Y)

	// Reset y
	b.yMatrix.Set(0, 0, newCenterX)
	b.yMatrix.Set(1, 0, newCenterY)

	// Reset u
	b.uMatrix.Set(0, 0, 0.0)
	b.uMatrix.Set(1, 0, 0.0)
	b.uMatrix.Set(2, 0, 0.0)
	b.uMatrix.Set(3, 0, 0.0)

	// Evaluate state
	state, err := b.pointTracker.Process(b.uMatrix, b.yMatrix)
	if err != nil {
		return errors.Wrap(err, "Can't process linear Kalman filter")
	}
	kalmanX, kalmanY := int(state.At(0, 0)), int(state.At(1, 0))
	b.CurrentRect = newbCast.CurrentRect
	b.Center = image.Point{kalmanX, kalmanY}
	diffX, diffY := kalmanX-newbCast.Center.X, kalmanY-newbCast.Center.Y
	b.CurrentRect = image.Rect(newbCast.CurrentRect.Min.X-diffX, newbCast.CurrentRect.Min.Y-diffY, newbCast.CurrentRect.Max.X-diffX, newbCast.CurrentRect.Max.Y-diffY)
	b.Area = newbCast.Area
	b.Diagonal = newbCast.Diagonal
	b.AspectRatio = newbCast.AspectRatio
	b.isStillBeingTracked = true
	b.isExists = true
	// Append new point to track
	b.Track = append(b.Track, b.Center)
	b.TrackTime = append(b.TrackTime, newbCast.TrackTime[len(newbCast.TrackTime)-1])
	// Restrict number of points in track (shift to the left)
	if len(b.Track) > b.maxPointsInTrack {
		b.Track = b.Track[1:]
	}
	return nil
}

func (sb *KalmanBlobie) GetID() uuid.UUID {
	return sb.ID
}

func (sb *KalmanBlobie) GetCenter() image.Point {
	return sb.Center
}

func (sb *KalmanBlobie) GetCurrentRect() image.Rectangle {
	return sb.CurrentRect
}

func (sb *KalmanBlobie) GetTrack() []image.Point {
	return sb.Track
}

func (sb *KalmanBlobie) GetTimestamps() []time.Time {
	return sb.TrackTime
}

func (sb *KalmanBlobie) GetDiagonal() float64 {
	return sb.Diagonal
}

func (sb *KalmanBlobie) GetPredictedNextPosition() image.Point {
	return sb.PredictedNextPosition
}

func (sb *KalmanBlobie) NoMatchTimes() int {
	return sb.noMatchTimes
}

func (sb *KalmanBlobie) Exists() bool {
	return sb.isExists
}

func (sb *KalmanBlobie) SetID(id uuid.UUID) {
	sb.ID = id
}

func (sb *KalmanBlobie) SetTracking(isStillBeingTracked bool) {
	sb.isStillBeingTracked = isStillBeingTracked
}

func (sb *KalmanBlobie) IncrementNoMatchTimes() {
	sb.noMatchTimes++
}

func (sb *KalmanBlobie) SetExists(isExists bool) {
	sb.isExists = isExists
}

// GetClassID Returns class identifier [KalmanBlobie]
func (b *KalmanBlobie) GetClassID() int {
	return b.classID
}

// GetClassName Returns class name [KalmanBlobie]
func (b *KalmanBlobie) GetClassName() string {
	return b.className
}

// SetDraw Sets options for drawing [KalmanBlobie]
func (b *KalmanBlobie) SetDraw(drawOptions *DrawOptions) {
	b.drawingOptions = drawOptions
}

// DrawTrack Draws blob's track [KalmanBlobie]
func (b *KalmanBlobie) DrawTrack(mat *gocv.Mat, optionalText ...string) {
	if b.drawingOptions == nil {
		b.drawingOptions = NewDrawOptionsDefault()
	}
	gocv.Rectangle(mat, b.CurrentRect, b.drawingOptions.BBoxColor.Color, b.drawingOptions.BBoxColor.Thickness)
	if b.isStillBeingTracked {
		for i := range b.Track {
			gocv.Circle(mat, b.Track[i], b.drawingOptions.CentroidColor.Radius, b.drawingOptions.CentroidColor.Color, b.drawingOptions.CentroidColor.Thickness)
		}
		shiftTextY := 0
		for i := len(optionalText) - 1; i >= 0; i-- {
			text := optionalText[i]
			if text != "" {
				anchor := image.Pt(b.CurrentRect.Min.X, b.CurrentRect.Min.Y-i*shiftTextY)
				textSize := gocv.GetTextSize(text, b.drawingOptions.TextColor.Font, b.drawingOptions.TextColor.Scale, b.drawingOptions.TextColor.Thickness)
				textRect := image.Rectangle{Min: image.Point{X: anchor.X, Y: anchor.Y - textSize.Y}, Max: image.Point{X: anchor.X + textSize.X, Y: anchor.Y}}
				gocv.Rectangle(mat, textRect, b.drawingOptions.BBoxColor.Color, b.drawingOptions.BBoxColor.Thickness)
				gocv.PutText(mat, text, anchor, b.drawingOptions.TextColor.Font, b.drawingOptions.TextColor.Scale, b.drawingOptions.TextColor.Color, b.drawingOptions.TextColor.Thickness)
				shiftTextY = textSize.Y
			}
		}
	}
}
