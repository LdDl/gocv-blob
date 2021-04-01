package blob

import (
	"image"
	"time"

	uuid "github.com/satori/go.uuid"
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

	// For array tracker
	drawingOptions *DrawOptions
	crossedLine    bool
}
