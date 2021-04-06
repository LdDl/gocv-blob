package blob

import (
	"time"
)

// BlobOptions Options for blob
type BlobOptions struct {
	ClassID          int
	ClassName        string
	MaxPointsInTrack int
	Time             time.Time
	TimeDeltaSeconds float64
}
