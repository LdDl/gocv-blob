package blob

import (
	"image"
	"image/color"

	uuid "github.com/satori/go.uuid"
	"gocv.io/x/gocv"
)

// Blobies - Array of blobs
type Blobies struct {
	Objects              map[uuid.UUID]*Blobie
	maxNoMatch           int
	minThresholdDistance float64
	maxPointsInTrack     int
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
	}
}

var arrayCounter = 1

// MatchToExisting Check if blob already exists
func (bt *Blobies) MatchToExisting(rects []image.Rectangle) {
	bt.prepare()
	blobies := make([]*Blobie, len(rects))
	for i := range blobies {
		blobies[i] = NewBlobie(rects[i], bt.maxPointsInTrack)
	}

	//log.Println("range before", len(*blobies), len(*bt))
	for i := range blobies {
		var minIndex uuid.UUID
		var minDistance = 200000.0
		for j := range (*bt).Objects {
			if bt.Objects[j].isStillBeingTracked == true {
				dist := distanceBetweenPoints(blobies[i].GetLastPoint(), (*bt).Objects[j].PredictedNextPosition)
				if dist < minDistance {
					minDistance = dist
					minIndex = j
				}
			}
		}
		// if minIndex == 0 {
		// 	log.Println("zero dawn", 0)
		// }
		// log.Println("dist", minDistance)
		if minDistance < blobies[i].Diagonal*0.5 {
			//log.Println("min dist", minDistance, (*b).Diagonal)
			bt.Objects[minIndex].Update(*blobies[i])
		} else {
			//bt.AddNew(blobies[i])
			bt.Register(blobies[i])
		}
	}
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
	newUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	bt.Objects[newUUID] = b
	return nil
}

// deregister - deregister blob with provided uuid
func (bt *Blobies) deregister(guid uuid.UUID) {
	delete(bt.Objects, guid)
}

// DrawTrack - Draw blob's track
func (b *Blobie) DrawTrack(mat *gocv.Mat, optionalText string) {
	if (*b).isStillBeingTracked {
		for i := range (*b).Track {
			gocv.Circle(mat, (*b).Track[i], 4, color.RGBA{255, 0, 0, 0}, 1)
		}
		gocv.Rectangle(mat, (*b).CurrentRect, color.RGBA{255, 255, 0, 0}, 2)
		pt := image.Pt((*b).CurrentRect.Min.X, (*b).CurrentRect.Min.Y)
		gocv.PutText(mat, optionalText, pt, gocv.FontHersheyPlain, 1.2, color.RGBA{0, 255, 0, 0}, 2)
	}
}

// IsCrossedTheLine - Check if blob crossed the HORIZONTAL line
func (b *Blobie) IsCrossedTheLine(horizontal int, direction bool) bool {
	trackLen := len(b.Track)
	if b.isStillBeingTracked == true && trackLen >= 2 && b.crossedLine == false {
		prevFrame := trackLen - 2
		currFrame := trackLen - 1
		if direction {
			if b.Track[prevFrame].Y <= horizontal && b.Track[currFrame].Y > horizontal { // TO us
				b.crossedLine = true
				return true
			}
		} else {
			if b.Track[prevFrame].Y > horizontal && b.Track[currFrame].Y <= horizontal { // FROM us
				b.crossedLine = true
				return true
			}
		}
	}
	return false
}
