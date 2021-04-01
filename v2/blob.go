package blob

import (
	"image"

	uuid "github.com/satori/go.uuid"
)

type Blob interface {
	GetCenter() image.Point
	GetPredictedNextPosition() image.Point
	GetDiagonal() float64
	Exists() bool
	NoMatchTimes() int
	IncrementNoMatchTimes()
	SetExists(isExists bool)
	SetTracking(isStillBeingTracked bool)
	SetID(id uuid.UUID)
	PredictNextPosition(n int)
	Update(newb Blob) error
}
