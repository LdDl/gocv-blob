package blob

import (
	"image"

	uuid "github.com/satori/go.uuid"
)

type Blobie interface {
	GetCenter() image.Point
	GetCurrentRect() image.Rectange
	GetPredictedNextPosition() image.Point
	GetDiagonal() float64
	GetClassID() int
	GetClassName() string
	Exists() bool
	NoMatchTimes() int
	IncrementNoMatchTimes()
	SetExists(isExists bool)
	SetTracking(isStillBeingTracked bool)
	SetID(id uuid.UUID)
	PredictNextPosition(n int)
	Update(newb Blobie) error
	SetDraw(drawOptions *DrawOptions)
	DrawTrack(mat *gocv.Mat, optionalText string)
}
