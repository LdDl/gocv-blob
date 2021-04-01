package blob

import (
	"image/color"

	"gocv.io/x/gocv"
)

// DrawOptions Options for blob.DrawTrack function
type DrawOptions struct {
	BBoxColor     DrawBBoxOptions
	CentroidColor DrawCentroidOptions
	TextColor     DrawTextOptions
}

// DrawBBoxOptions Options for bounding box of blob
type DrawBBoxOptions struct {
	Color     color.RGBA
	Thickness int
}

// DrawCentroidOptions Options for centroid of blob
type DrawCentroidOptions struct {
	Color     color.RGBA
	Radius    int
	Thickness int
}

// DrawTextOptions Options for text if top left corner bounding box of blob
type DrawTextOptions struct {
	Color     color.RGBA
	Scale     float64
	Thickness int
	Font      gocv.HersheyFont
}

// NewDrawOptionsDefault Returns pointer to new DrawOptions (default)
func NewDrawOptionsDefault() *DrawOptions {
	bbox := DrawBBoxOptions{
		Color:     color.RGBA{255, 255, 0, 0},
		Thickness: 2,
	}
	centroid := DrawCentroidOptions{
		Color:     color.RGBA{255, 0, 0, 0},
		Radius:    4,
		Thickness: 2,
	}
	text := DrawTextOptions{
		Color:     color.RGBA{0, 255, 0, 0},
		Scale:     1.2,
		Thickness: 2,
		Font:      gocv.FontHersheyPlain,
	}
	opts := DrawOptions{
		BBoxColor:     bbox,
		CentroidColor: centroid,
		TextColor:     text,
	}
	return &opts
}

// NewDrawOptions Returns pointer to new DrawOptions (custom)
func NewDrawOptions(bbox DrawBBoxOptions, centroid DrawCentroidOptions, text DrawTextOptions) *DrawOptions {
	opts := DrawOptions{
		BBoxColor:     bbox,
		CentroidColor: centroid,
		TextColor:     text,
	}
	return &opts
}
