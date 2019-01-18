package blob

import (
	"image"

	uuid "github.com/satori/go.uuid"
)

// CentroidTracker - See ref. https://www.pyimagesearch.com/2018/07/23/simple-object-tracking-with-opencv/
type CentroidTracker struct {
	Objects              map[uuid.UUID]*Blobie
	maxNoMatch           int
	minThresholdDistance float64
	maxPointsInTrack     int
}

// NewCentroidTracker - Constructor for CentroidTracker
func NewCentroidTracker(maxNoMatch int, minThresholdDistance float64, maxPointsInTrack int) *CentroidTracker {
	return &CentroidTracker{
		Objects:              make(map[uuid.UUID]*Blobie),
		maxNoMatch:           maxNoMatch,
		minThresholdDistance: minThresholdDistance,
		maxPointsInTrack:     maxPointsInTrack,
	}
}

// NewCentroidTrackerDefaults - Constructor for CentroidTracker (default values)
//
// Default values are:
// maxNoMatch = 30
// minThresholdDistance = 15
// maxPointsInTrack = 10
//
func NewCentroidTrackerDefaults() *CentroidTracker {
	return &CentroidTracker{
		Objects:              make(map[uuid.UUID]*Blobie),
		maxNoMatch:           30,
		minThresholdDistance: 15,
		maxPointsInTrack:     10,
	}
}

// Register - Register new blob
func (ct *CentroidTracker) Register(b ...*Blobie) error {
	for i := range b {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		ct.Objects[newUUID] = b[i]
	}
	return nil
}

// deregister - deregister blob with provided uuid
func (ct *CentroidTracker) deregister(guid uuid.UUID) {
	delete(ct.Objects, guid)
}

// Update - Update blobs in tracker. See ref. https://www.pyimagesearch.com/2018/07/23/simple-object-tracking-with-opencv/
func (ct *CentroidTracker) Update(rects []image.Rectangle) map[uuid.UUID]*Blobie {
	if len(rects) == 0 {
		ct.allDisappear()
		return ct.Objects
	}

	blobies := make([]*Blobie, len(rects))
	for i := range blobies {
		blobies[i] = NewBlobie(rects[i], ct.maxPointsInTrack)
		blobies[i].isExists = false
	}

	if len(ct.Objects) == 0 {
		ct.Register(blobies...)
		return ct.Objects
	}

	for i := range ct.Objects {
		minDist := ct.minThresholdDistance
		minIdx := -1
		for j := range blobies {
			dist := distanceBetweenBlobies(ct.Objects[i], blobies[j])
			if minDist > dist {
				minIdx = j
				minDist = dist
			}
		}
		if minIdx == -1 {
			ct.Objects[i].noMatchTimes++
			if ct.Objects[i].noMatchTimes > ct.maxNoMatch {
				ct.deregister(i)
			}
			continue
		}
		blobies[minIdx].isExists = true
		ct.Objects[i] = blobies[minIdx]
		ct.Objects[i].noMatchTimes = 0
	}

	for i := range blobies {
		if blobies[i].isExists == false {
			ct.Register(blobies[i])
		}
	}

	return ct.Objects
}

func (ct *CentroidTracker) allDisappear() {
	for u := range ct.Objects {
		ct.Objects[u].noMatchTimes++
		if ct.Objects[u].noMatchTimes > ct.maxNoMatch {
			ct.deregister(u)
		}
	}
}
