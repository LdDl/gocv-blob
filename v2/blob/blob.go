package blob

import (
	"image"

	uuid "github.com/satori/go.uuid"
	"gocv.io/x/gocv"
)

type Blobie interface {
	GetID() uuid.UUID
	GetCenter() image.Point
	GetCurrentRect() image.Rectangle
	GetPredictedNextPosition() image.Point
	GetTrack() []image.Point
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
	IsCrossedTheLine(vertical, leftX, rightX int, direction bool) bool
	IsCrossedTheLineWithShift(vertical, leftX, rightX int, direction bool, shift int) bool
}
