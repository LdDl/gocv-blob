package blob

import (
	"image"
	"testing"
)

func TestKalmanPredictPos(t *testing.T) {
	var (
		testPoints = [][]int{
			[]int{0, 0},
			[]int{1, 1},
			[]int{2, 2},
			[]int{4, 4},
			[]int{6, 6},
			[]int{9, 9},
			[]int{11, 11},
			[]int{16, 16},
			[]int{20, 20},
		}
		kalmanFilteredPoints = [][]int{
			[]int{0, 0},
			[]int{0, 0},
			[]int{1, 1},
			[]int{3, 3},
			[]int{5, 5},
			[]int{7, 7},
			[]int{10, 10},
			[]int{13, 13},
			[]int{17, 17},
		}
		correctPredictions = [][]int{
			[]int{0, 0},
			[]int{0, 0},
			[]int{0, 0},
			[]int{1, 1},
			[]int{4, 4},
			[]int{6, 6},
			[]int{8, 8},
			[]int{12, 12},
			[]int{15, 15},
		}
	)

	maxPointsInTrack := 150
	classID := 1
	className := "just_an_object"
	maxNoMatch := 5

	rectHalfHeight := 30
	rectHalfWidth := 75

	commonOptions := BlobOptions{
		ClassID:          classID,
		ClassName:        className,
		MaxPointsInTrack: maxPointsInTrack,
		TimeDeltaSeconds: 1.0,
	}

	var b Blobie

	for i := range testPoints {
		centerOne := testPoints[i]
		rectOne := image.Rect(centerOne[0]-rectHalfWidth, centerOne[1]-rectHalfHeight, centerOne[0]+rectHalfWidth, centerOne[1]+rectHalfHeight)
		blobOne := NewKalmanBlobie(rectOne, &commonOptions)
		if b == nil {
			// Fill data on first iteration
			b = blobOne
		}
		b.PredictNextPosition(maxNoMatch)
		b.Update(blobOne)
		forCheck := b.(*KalmanBlobie)
		smoothedCenter := kalmanFilteredPoints[i]
		if forCheck.Center.X != smoothedCenter[0] {
			t.Errorf("Center.X on %d-th iteration should be %d, but got %d", i, smoothedCenter[0], forCheck.Center.X)
		}
		if forCheck.Center.Y != smoothedCenter[1] {
			t.Errorf("Center.Y on %d-th iteration should be %d, but got %d", i, smoothedCenter[1], forCheck.Center.Y)
		}
		if forCheck.PredictedNextPosition.X != correctPredictions[i][0] {
			t.Errorf("PredictedNextPosition.X on %d-th iteration should be %d, but got %d", i, correctPredictions[i][0], forCheck.PredictedNextPosition.X)
		}
		if forCheck.PredictedNextPosition.Y != correctPredictions[i][1] {
			t.Errorf("PredictedNextPosition.Y on %d-th iteration should be %d, but got %d", i, correctPredictions[i][1], forCheck.PredictedNextPosition.Y)
		}
	}
}
