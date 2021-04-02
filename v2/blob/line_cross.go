package blob

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
