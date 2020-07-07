package blob

import (
	"image"
	"math"

	uuid "github.com/satori/go.uuid"
)

// Blobies - Array of blobs
type Blobies struct {
	Objects              map[uuid.UUID]*Blobie
	maxNoMatch           int
	minThresholdDistance float64
	maxPointsInTrack     int

	DrawingOptions *DrawOptions
}

// NewBlobiesDefaults - Constructor for Blobies (default values)
//
// Default values are:
// maxNoMatch = 5
// minThresholdDistance = 15
// maxPointsInTrack = 10
//
func NewBlobiesDefaults() *Blobies {
	return &Blobies{
		Objects:              make(map[uuid.UUID]*Blobie),
		maxNoMatch:           5,
		minThresholdDistance: 15,
		maxPointsInTrack:     10,
		DrawingOptions:       NewDrawOptionsDefault(),
	}
}

// MatchToExisting Check if some of blobs already exists
func (bt *Blobies) MatchToExisting(rects []image.Rectangle) {
	bt.prepare()
	blobies := make([]*Blobie, len(rects))
	for i := range blobies {
		blobies[i] = NewBlobie(rects[i], bt.maxPointsInTrack)
		blobies[i].drawingOptions = bt.DrawingOptions
	}

	for i := range blobies {
		minUUID := uuid.UUID{}
		minDistance := math.MaxFloat64
		for j := range (*bt).Objects {
			dist := distanceBetweenPoints(blobies[i].Center, (*bt).Objects[j].Center)
			distPredicted := distanceBetweenPoints(blobies[i].Center, (*bt).Objects[j].PredictedNextPosition)
			dist = minf64(dist, distPredicted)
			if dist < minDistance {
				minDistance = dist
				minUUID = j
			}
		}
		if minDistance < blobies[i].Diagonal*0.5 || minDistance < bt.minThresholdDistance {
			bt.Objects[minUUID].Update(*blobies[i])
		} else {
			bt.Register(blobies[i])
		}
	}

	bt.RefreshNoMatch()
}

func (bt *Blobies) RefreshNoMatch() {
	for i, b := range (*bt).Objects {
		if b.isExists == false {
			b.noMatchTimes++
		}
		if b.noMatchTimes >= 5 {
			b.isStillBeingTracked = false
			bt.deregister(i)
		}
	}
}

func (bt *Blobies) prepare() {
	for i := range bt.Objects {
		bt.Objects[i].isExists = false
		bt.Objects[i].PredictNextPosition(bt.maxNoMatch)
	}
}

// Register - Register new blob
func (bt *Blobies) Register(b *Blobie) error {
	newUUID := uuid.NewV4()
	b.ID = newUUID
	bt.Objects[newUUID] = b
	return nil
}

// deregister - deregister blob with provided uuid
func (bt *Blobies) deregister(guid uuid.UUID) {
	delete(bt.Objects, guid)
}

// IsCrossedTheLine - Check if blob crossed the HORIZONTAL line
func (b *Blobie) IsCrossedTheLine(vertical, leftX, rightX int, direction bool) bool {
	trackLen := len(b.Track)
	if b.isStillBeingTracked == true && trackLen >= 2 && b.crossedLine == false {
		prevFrame := trackLen - 2
		currFrame := trackLen - 1
		if b.Track[currFrame].X >= leftX && b.Track[currFrame].X <= rightX {
			if direction {

				if b.Track[prevFrame].Y <= vertical && b.Track[currFrame].Y > vertical { // TO us
					b.crossedLine = true
					return true
				}
			} else {
				if b.Track[prevFrame].Y > vertical && b.Track[currFrame].Y <= vertical { // FROM us
					b.crossedLine = true
					return true
				}
			}
		}
	}
	return false
}

// IsCrossedTheLineWithShift - Check if blob crossed the HORIZONTAL line with shift along the Y-axis
// Purpose of this for "predicative" cropping when detection line very close to bottom of image
func (b *Blobie) IsCrossedTheLineWithShift(vertical, leftX, rightX int, direction bool, shift int) bool {
	trackLen := len(b.Track)
	if b.isStillBeingTracked == true && trackLen >= 2 && b.crossedLine == false {
		prevFrame := trackLen - 2
		currFrame := trackLen - 1
		if b.Track[currFrame].X >= leftX && b.Track[currFrame].X <= rightX {
			if direction {
				if (b.Track[prevFrame].Y+shift) <= vertical && (b.Track[currFrame].Y+shift) > vertical { // TO us
					b.crossedLine = true
					return true
				}
			} else {
				if (b.Track[prevFrame].Y+shift) > vertical && (b.Track[currFrame].Y+shift) <= vertical { // FROM us
					b.crossedLine = true
					return true
				}
			}
		}
	}
	return false
}
