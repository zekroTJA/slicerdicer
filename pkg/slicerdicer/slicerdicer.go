// Package slicerdicer provides a very simple API to
// crop image objects and slice image objects into
// same-sized pieces.
package slicerdicer

import (
	"errors"
	"image"
	"image/draw"
)

var (
	ErrPointOutOfBounds     = errors.New("point is out of bounds")
	ErrDimensionOutOfBounds = errors.New("dimension is out of bounds")
)

// Crop returns an image which is cropped by the given
// dimensions dx and dy starting at the point x and y.
//
// Errors are returned when the given point and dimensions
// are out of the image bounds.
func Crop(img image.Image, x, y, dx, dy int) (image.Image, error) {
	bounds := img.Bounds()

	if x < 0 || x > bounds.Dx() || y < 0 || y > bounds.Dy() {
		return nil, ErrPointOutOfBounds
	}

	dx = x + dx
	dy = y + dy

	if dx < 1 || dx > bounds.Dx() || dy < 1 || dy > bounds.Dy() {
		return nil, ErrDimensionOutOfBounds
	}

	rect := image.Rect(x, y, dx, dy)
	res := image.NewRGBA(rect)

	draw.Draw(res, rect, img, rect.Min, draw.Src)

	return res, nil
}

// Slice takes an image and an ammount of slices per side.
// The passed image is then sliced into partsPerSide * partsPerSide
// pieces with the size of imgX / partsPerSide times imgY / partsPerSide.
// The resulting pieces are returned in a two-dimensional array of
// image instances arranged in the order of the original image.
//
// Errors are returned if the calculated sizes are out of bounds
// of the passed image.
func Slice(img image.Image, partsPerSide int) ([][]image.Image, error) {
	bounds := img.Bounds()

	xLen := bounds.Dx() / partsPerSide
	yLen := bounds.Dy() / partsPerSide

	res := make([][]image.Image, partsPerSide)
	for i, _ := range res {
		res[i] = make([]image.Image, partsPerSide)
	}

	for y := 0; y < partsPerSide; y++ {
		for x := 0; x < partsPerSide; x++ {
			r, err := Crop(img, x*xLen, y*yLen, xLen, yLen)
			if err != nil {
				return nil, err
			}
			res[y][x] = r
		}
	}

	return res, nil
}
