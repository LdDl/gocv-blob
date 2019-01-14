package blob

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

// Blobie - Detected candidate to be target object
type Blobie struct {
	CurrentRect                         image.Rectangle
	CurrentContour                      []image.Point
	Track                               []image.Point
	PredictedNextPosition               image.Point
	Diagonal                            float64
	AspectRatio                         float64
	Area                                float64
	IsExists                            bool
	IsStillBeingTracked                 bool
	Counted                             bool
	NumOfConsecutiveFramesWithoutAMatch int
}

// Blobies - Array of blobs
type Blobies []*Blobie

// NewBlobieFromRect - Blobie constructor. Rect as provider
func NewBlobieFromRect(rect *image.Rectangle) *Blobie {
	var currentCenter image.Point
	var rectWidth = (*rect).Dx()
	var rectHeight = (*rect).Dy()
	currentCenter.X = ((*rect).Min.X*2 + rectWidth) / 2
	currentCenter.Y = ((*rect).Min.Y*2 + rectHeight) / 2
	var b = Blobie{
		CurrentRect:                         (*rect),
		Track:                               []image.Point{currentCenter},
		Area:                                float64(rectWidth * rectHeight),
		Diagonal:                            math.Sqrt(math.Pow(float64(rectWidth), 2) + math.Pow(float64(rectHeight), 2)),
		AspectRatio:                         float64(rectWidth) / float64(rectHeight),
		IsStillBeingTracked:                 true,
		Counted:                             false,
		IsExists:                            true,
		NumOfConsecutiveFramesWithoutAMatch: 0,
	}
	return &b
}

// NewBlobieFromContour - Blobie constructor. Contour (pointer to array of image.Point) as provider
func NewBlobieFromContour(contour *[]image.Point) *Blobie {
	var b Blobie
	// gocv.ContourArea((*contour))
	// @todo !!!
	return &b
}

// IsCrossedTheLine - Check if blob crossed the HORIZONTAL line
func (b *Blobie) IsCrossedTheLine(horizontal int, counter *int, direction bool) bool {
	if (*b).IsStillBeingTracked == true && len((*b).Track) >= 2 && (*b).Counted == false {
		prevFrame := len((*b).Track) - 2
		currFrame := len((*b).Track) - 1
		if direction {
			if (*b).Track[prevFrame].Y <= horizontal && (*b).Track[currFrame].Y > horizontal { // TO us
				*counter++
				b.AsCounted()
				return true
			}
		} else {
			if (*b).Track[prevFrame].Y > horizontal && (*b).Track[currFrame].Y <= horizontal { // FROM us
				*counter++
				b.AsCounted()
				return true
			}
		}
	}
	return false
}

// DrawTrack - Draw blob's track
func (b *Blobie) DrawTrack(mat *gocv.Mat, optionalText string) {
	if (*b).IsStillBeingTracked == true {
		for i := range (*b).Track {
			gocv.Circle(mat, (*b).Track[i], 4, color.RGBA{255, 0, 0, 0}, 1)
		}
		gocv.Rectangle(mat, (*b).CurrentRect, color.RGBA{0, 255, 255, 0}, 2)
		pt := image.Pt((*b).CurrentRect.Min.X, (*b).CurrentRect.Min.Y)
		gocv.PutText(mat, optionalText, pt, gocv.FontHersheyPlain, 1.2, color.RGBA{0, 255, 0, 0}, 2)
	}
}

// GetRect - Get blob's rect
func (b *Blobie) GetRect(mat *gocv.Mat, id string) image.Rectangle {
	if (*b).IsStillBeingTracked == true {
		return (*b).CurrentRect
	}
	return image.Rectangle{}
}

// StopTracking - Stop tracking object
func (b *Blobie) StopTracking() {
	(*b).IsStillBeingTracked = false
}

// AsCounted - Set object as counted
func (b *Blobie) AsCounted() {
	(*b).Counted = true
}

// PredictNextPosition - Predict next position
func (b *Blobie) PredictNextPosition() {
	account := min(5, int64(len((*b).Track)))
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

// MatchToExisting Check if blob already exists
func (bExisting *Blobies) MatchToExisting(bCurrent *Blobies) {
	for _, b := range *bExisting {
		(*b).IsExists = false
		(*b).PredictNextPosition()
	}
	for _, b := range *bCurrent {
		var intIndexOfLeastDistance = 0
		var dblLeastDistance = 200000.0
		for i := range *bExisting {
			if (*bExisting)[i].IsStillBeingTracked == true {
				dblDistance := distanceBetweenPoints((b).Track[len((*b).Track)-1], (*bExisting)[i].PredictedNextPosition)
				if dblDistance < dblLeastDistance {
					dblLeastDistance = dblDistance
					intIndexOfLeastDistance = i
				}
			}
		}
		if dblLeastDistance < (*b).Diagonal*0.5 {
			(*bExisting).AddToExisting(b, intIndexOfLeastDistance)
		} else {
			(*bExisting).AddNew(b)
		}
	}
	for _, b := range *bExisting {
		if (*b).IsExists == false {
			(*b).NumOfConsecutiveFramesWithoutAMatch++
		}
		if (*b).NumOfConsecutiveFramesWithoutAMatch >= 5 {
			(*b).IsStillBeingTracked = false
		}
	}

}

// AddToExisting Add blob to existing blobs
func (bExisting *Blobies) AddToExisting(bCurrent *Blobie, intIndex int) {
	(*bExisting)[intIndex].CurrentRect = (*bCurrent).CurrentRect
	(*bExisting)[intIndex].Track = append((*bExisting)[intIndex].Track, (*bCurrent).Track[len((*bCurrent).Track)-1])
	(*bExisting)[intIndex].Diagonal = (*bCurrent).Diagonal
	(*bExisting)[intIndex].AspectRatio = (*bCurrent).AspectRatio
	(*bExisting)[intIndex].IsStillBeingTracked = true
	(*bExisting)[intIndex].IsExists = true
}

// AddNew Add new blob to array of blobs
func (bExisting *Blobies) AddNew(bCurrent *Blobie) {
	(*bCurrent).IsExists = true
	(*bExisting) = append((*bExisting), bCurrent)
}

// Flush - Delete all blobies
func (bExisting *Blobies) Flush() {
	(*bExisting) = Blobies{}
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func minf(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

func maxf(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

func distanceBetweenPoints(p1 image.Point, p2 image.Point) float64 {
	intX := math.Abs(float64(p1.X - p2.X))
	intY := math.Abs(float64(p1.Y - p2.Y))
	return math.Sqrt(math.Pow(intX, 2) + math.Pow(intY, 2))
}
