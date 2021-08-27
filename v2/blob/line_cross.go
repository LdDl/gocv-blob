package blob

type PointsOrientation int

const (
	Collinear = iota
	Clockwise
	CounterClockwise
)

func maxInt(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// isOnSegment Checks if point Q lies on segment PR
// Input: three colinear points Q, Q and R
func isOnSegment(Px, Py, Qx, Qy, Rx, Ry int) bool {
	if Qx <= maxInt(Px, Rx) && Qx >= minInt(Px, Rx) && Qy <= maxInt(Py, Ry) && Qy >= minInt(Py, Ry) {
		return true
	}
	return false
}

// getOrientation Gets orientations of points P -> Q -> R.
// Possible output values: Collinear / Clockwise or CounterClockwise
// Input: points P, Q and R in provided order
func getOrientation(Px, Py, Qx, Qy, Rx, Ry int) PointsOrientation {
	val := (Qy-Py)*(Rx-Qx) - (Qx-Px)*(Ry-Qy)
	if val == 0 {
		return Collinear
	}
	if val > 0 {
		return Clockwise
	}
	return CounterClockwise // if it's neither collinear nor clockwise
}

// isIntersects Checks if segments intersect each other
// Input:
// firstPx, firstPy, firstQx, firstQy === first segment
// secondPx, secondPy, secondQx, secondQy === second segment
/*
Notation
	P1 = (firstPx, firstPy)
	Q1 = (firstQx, firstQy)
	P2 = (secondPx, secondPy)
	Q2 = (secondQx, secondQy)
*/
func isIntersects(firstPx, firstPy, firstQx, firstQy, secondPx, secondPy, secondQx, secondQy int) bool {
	// Find the four orientations needed for general case and special ones
	o1 := getOrientation(firstPx, firstPy, firstQx, firstQy, secondPx, secondPy)
	o2 := getOrientation(firstPx, firstPy, firstQx, firstQy, secondQx, secondQy)
	o3 := getOrientation(secondPx, secondPy, secondQx, secondQy, firstPx, firstPy)
	o4 := getOrientation(secondPx, secondPy, secondQx, secondQy, firstQx, firstQy)

	// General case
	if o1 != o2 && o3 != o4 {
		return true
	}

	/* Special cases */
	// P1, Q1, P2 are colinear and P2 lies on segment P1-Q1
	if o1 == Collinear && isOnSegment(firstPx, firstPy, secondPx, secondPy, firstQx, firstQy) {
		return true
	}
	// P1, Q1 and Q2 are colinear and Q2 lies on segment P1-Q1
	if o2 == Collinear && isOnSegment(firstPx, firstPy, secondQx, secondQy, firstQx, firstQy) {
		return true
	}
	// P2, Q2 and P1 are colinear and P1 lies on segment P2-Q2
	if o3 == Collinear && isOnSegment(secondPx, secondPy, firstPx, firstPy, secondQx, secondQy) {
		return true
	}
	// P2, Q2 and Q1 are colinear and Q1 lies on segment P2-Q2
	if o4 == Collinear && isOnSegment(secondPx, secondPy, firstQx, firstQy, secondQx, secondQy) {
		return true
	}
	// Segments do not intersect
	return false
}

// IsCrossedTheLine - Check if blob crossed the HORIZONTAL line
func (b *SimpleBlobie) IsCrossedTheLine(vertical, leftX, rightX int, direction bool) bool {
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
func (b *SimpleBlobie) IsCrossedTheLineWithShift(vertical, leftX, rightX int, direction bool, shift int) bool {
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

// IsCrossedTheObliqueLine - Check if blob crossed the OBLIQUE line
// This should be used when lineStart.Y != lineEnd.Y
func (b *SimpleBlobie) IsCrossedTheObliqueLine(leftX, leftY, rightX, rightY int, direction bool) bool {
	trackLen := len(b.Track)
	if b.isStillBeingTracked == true && trackLen >= 2 && b.crossedLine == false {
		prevFrame := trackLen - 2
		currFrame := trackLen - 1
		// First segment is: P1 = (b.Track[prevFrame].X, b.Track[prevFrame].Y), Q1 = (b.Track[currFrame].X, b.Track[currFrame].Y)
		// Second segment is: P2 = (leftX, leftY), Q2 = (rightX, rightY)
		if isIntersects(b.Track[prevFrame].X, b.Track[prevFrame].Y, b.Track[currFrame].X, b.Track[currFrame].Y, leftX, leftY, rightX, rightY) {
			if direction {
				if b.Track[currFrame].Y > b.Track[prevFrame].Y { // TO us
					b.crossedLine = true
					return true
				}
			} else {
				if b.Track[currFrame].Y <= b.Track[prevFrame].Y { // FROM us
					b.crossedLine = true
					return true
				}
			}
		}
	}
	return false
}

// IsCrossedTheObliqueLine - Check if blob crossed the OBLIQUE line with shift along the Y-axis
// This should be used when lineStart.Y != lineEnd.Y
// Purpose of shifting: for "predicative" cropping when detection line very close to bottom of image
func (b *SimpleBlobie) IsCrossedTheObliqueLineWithShift(leftX, leftY, rightX, rightY int, direction bool, shift int) bool {
	trackLen := len(b.Track)
	if b.isStillBeingTracked == true && trackLen >= 2 && b.crossedLine == false {
		prevFrame := trackLen - 2
		currFrame := trackLen - 1
		// First segment is: P1 = (b.Track[prevFrame].X, b.Track[prevFrame].Y + shift), Q1 = (b.Track[currFrame].X, b.Track[currFrame].Y + shift)
		// Second segment is: P2 = (leftX, leftY), Q2 = (rightX, rightY)
		if isIntersects(b.Track[prevFrame].X, b.Track[prevFrame].Y+shift, b.Track[currFrame].X, b.Track[currFrame].Y+shift, leftX, leftY, rightX, rightY) {
			if direction {
				if b.Track[currFrame].Y > b.Track[prevFrame].Y { // TO us
					b.crossedLine = true
					return true
				}
			} else {
				if b.Track[currFrame].Y <= b.Track[prevFrame].Y { // FROM us
					b.crossedLine = true
					return true
				}
			}
		}
	}
	return false
}

// IsCrossedTheLine - Check if blob crossed the HORIZONTAL line
func (b *KalmanBlobie) IsCrossedTheLine(vertical, leftX, rightX int, direction bool) bool {
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
func (b *KalmanBlobie) IsCrossedTheLineWithShift(vertical, leftX, rightX int, direction bool, shift int) bool {
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

// IsCrossedTheObliqueLine - Check if blob crossed the OBLIQUE line
// This should be used when lineStart.Y != lineEnd.Y
func (b *KalmanBlobie) IsCrossedTheObliqueLine(leftX, leftY, rightX, rightY int, direction bool) bool {
	trackLen := len(b.Track)
	if b.isStillBeingTracked == true && trackLen >= 2 && b.crossedLine == false {
		prevFrame := trackLen - 2
		currFrame := trackLen - 1
		// First segment is: P1 = (b.Track[prevFrame].X, b.Track[prevFrame].Y), Q1 = (b.Track[currFrame].X, b.Track[currFrame].Y)
		// Second segment is: P2 = (leftX, leftY), Q2 = (rightX, rightY)
		if isIntersects(b.Track[prevFrame].X, b.Track[prevFrame].Y, b.Track[currFrame].X, b.Track[currFrame].Y, leftX, leftY, rightX, rightY) {

			if direction {
				if b.Track[currFrame].Y > b.Track[prevFrame].Y { // TO us
					b.crossedLine = true
					return true
				}
			} else {
				if b.Track[currFrame].Y <= b.Track[prevFrame].Y { // FROM us
					b.crossedLine = true
					return true
				}
			}
		}
	}
	return false
}

// IsCrossedTheObliqueLine - Check if blob crossed the OBLIQUE line with shift along the Y-axis
// This should be used when lineStart.Y != lineEnd.Y
// Purpose of shifting: for "predicative" cropping when detection line very close to bottom of image
func (b *KalmanBlobie) IsCrossedTheObliqueLineWithShift(leftX, leftY, rightX, rightY int, direction bool, shift int) bool {
	trackLen := len(b.Track)
	if b.isStillBeingTracked == true && trackLen >= 2 && b.crossedLine == false {
		prevFrame := trackLen - 2
		currFrame := trackLen - 1
		// First segment is: P1 = (b.Track[prevFrame].X, b.Track[prevFrame].Y + shift), Q1 = (b.Track[currFrame].X, b.Track[currFrame].Y + shift)
		// Second segment is: P2 = (leftX, leftY), Q2 = (rightX, rightY)
		if isIntersects(b.Track[prevFrame].X, b.Track[prevFrame].Y+shift, b.Track[currFrame].X, b.Track[currFrame].Y+shift, leftX, leftY, rightX, rightY) {
			if direction {
				if b.Track[currFrame].Y > b.Track[prevFrame].Y { // TO us
					b.crossedLine = true
					return true
				}
			} else {
				if b.Track[currFrame].Y <= b.Track[prevFrame].Y { // FROM us
					b.crossedLine = true
					return true
				}
			}
		}
	}
	return false
}
