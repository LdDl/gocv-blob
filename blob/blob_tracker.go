package blob

import (
	"image"

	uuid "github.com/satori/go.uuid"
)

// BlobiesTracker - Blob tracker
type BlobiesTracker struct {
	Objects              map[uuid.UUID]*Blobie
	maxNoMatch           int
	minThresholdDistance float64
	maxPointsInTrack     int
}

// NewBlobiesTracker - Constructor for BlobiesTracker
func NewBlobiesTracker(maxNoMatch int, minThresholdDistance float64, maxPointsInTrack int) *BlobiesTracker {
	return &BlobiesTracker{
		Objects:              make(map[uuid.UUID]*Blobie),
		maxNoMatch:           maxNoMatch,
		minThresholdDistance: minThresholdDistance,
		maxPointsInTrack:     maxPointsInTrack,
	}
}

// NewBlobiesTrackerDefaults - Constructor for BlobiesTracker (default values)
//
// Default values are:
// maxNoMatch = 5
// minThresholdDistance = 15
// maxPointsInTrack = 10
//
func NewBlobiesTrackerDefaults() *BlobiesTracker {
	return &BlobiesTracker{
		Objects:              make(map[uuid.UUID]*Blobie),
		maxNoMatch:           5,
		minThresholdDistance: 15,
		maxPointsInTrack:     10,
	}
}

// Register - Register new blob
func (bt *BlobiesTracker) Register(b ...*Blobie) error {
	for i := range b {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		bt.Objects[newUUID] = b[i]
	}
	return nil
}

// deregister - deregister blob with provided uuid
func (bt *BlobiesTracker) deregister(guid uuid.UUID) {
	delete(bt.Objects, guid)
}

// Update - Update blobs in tracker
func (bt *BlobiesTracker) Update(rects []image.Rectangle) map[uuid.UUID]*Blobie {

	// if len(rects) == 0 {
	// 	bt.allDisappear()
	// 	return bt.Objects
	// }

	bt.prepare()
	blobies := make([]*Blobie, len(rects))
	for i := range blobies {
		blobies[i] = NewBlobie(rects[i], bt.maxPointsInTrack)
	}

	// if len(bt.Objects) == 0 {
	// 	bt.Register(blobies...)
	// 	return bt.Objects
	// }

	//log.Println("range before panic", len(blobies), len(bt.Objects))
	for i := range blobies {
		var minIndex uuid.UUID
		var minDistance = 200000.0
		for j := range bt.Objects {
			if bt.Objects[j].isStillBeingTracked {
				dist := distanceBetweenPoints(bt.Objects[j].GetLastPoint(), bt.Objects[j].PredictedNextPosition)
				if dist < minDistance {
					minDistance = dist
					minIndex = j
				}
			}
		}
		// log.Println("dist panic", minDistance)
		if minDistance < blobies[i].Diagonal*0.5 {
			//log.Println("min dist panic", minDistance, blobies[i].Diagonal)
			bt.Objects[minIndex].Update(*blobies[i])
		} else {
			// log.Println(minDistance, blobies[i].Diagonal*0.5)
			bt.Register(blobies[i])
		}
	}
	// log.Println("range after", len(blobies), len(bt.Objects))
	for i := range bt.Objects {
		if bt.Objects[i].isExists == false {
			bt.Objects[i].noMatchTimes++
		}
		if bt.Objects[i].noMatchTimes > bt.maxNoMatch {
			bt.Objects[i].isStillBeingTracked = false
			bt.deregister(i)
		}
	}

	return bt.Objects
}

func (bt *BlobiesTracker) prepare() {
	for i := range bt.Objects {
		bt.Objects[i].isExists = false
		bt.Objects[i].PredictNextPosition(bt.maxNoMatch)
	}
}

func (bt *BlobiesTracker) allDisappear() {
	for u := range bt.Objects {
		bt.Objects[u].noMatchTimes++
		if bt.Objects[u].noMatchTimes > bt.maxNoMatch {
			bt.deregister(u)
		}
	}
}
