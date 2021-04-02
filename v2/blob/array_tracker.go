package blob

import (
	"math"

	uuid "github.com/satori/go.uuid"
)

// Blobies - Array of blobs
type Blobies struct {
	Objects              map[uuid.UUID]Blobie
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
		Objects:              make(map[uuid.UUID]Blobie),
		maxNoMatch:           5,
		minThresholdDistance: 15,
		maxPointsInTrack:     10,
		DrawingOptions:       NewDrawOptionsDefault(),
	}
}

// MatchToExisting Check if some of blobs already exists
func (bt *Blobies) MatchToExisting(blobies []Blobie) {
	bt.prepare()
	for i := range blobies {
		minUUID := uuid.UUID{}
		minDistance := math.MaxFloat64
		for j := range (*bt).Objects {
			dist := distanceBetweenPoints(blobies[i].GetCenter(), (*bt).Objects[j].GetCenter())
			distPredicted := distanceBetweenPoints(blobies[i].GetCenter(), (*bt).Objects[j].GetPredictedNextPosition())
			dist = minf64(dist, distPredicted)
			if dist < minDistance {
				minDistance = dist
				minUUID = j
			}
		}
		if minDistance < blobies[i].GetDiagonal()*0.5 || minDistance < bt.minThresholdDistance {
			bt.Objects[minUUID].Update(blobies[i])
		} else {
			bt.Register(blobies[i])
		}
	}
	bt.RefreshNoMatch()
}

// RefreshNoMatch - Refresh state of each blob
func (bt *Blobies) RefreshNoMatch() {
	for i, b := range (*bt).Objects {
		if b.Exists() == false {
			b.IncrementNoMatchTimes()
		}
		if b.NoMatchTimes() >= 5 {
			b.SetTracking(false)
			bt.deregister(i)
		}
	}
}

func (bt *Blobies) prepare() {
	for i := range bt.Objects {
		bt.Objects[i].SetExists(false)
		bt.Objects[i].PredictNextPosition(bt.maxNoMatch)
	}
}

// Register - Register new blob
func (bt *Blobies) Register(b Blobie) error {
	newUUID := uuid.NewV4()
	b.SetID(newUUID)
	bt.Objects[newUUID] = b
	return nil
}

// deregister - deregister blob with provided uuid
func (bt *Blobies) deregister(guid uuid.UUID) {
	delete(bt.Objects, guid)
}
