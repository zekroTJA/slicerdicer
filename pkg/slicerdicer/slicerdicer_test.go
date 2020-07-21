package slicerdicer

import (
	"image"
	"testing"
)

var (
	testRect  = image.Rect(0, 0, 1000, 1000)
	testImage = image.NewRGBA(testRect)
)

func assert(t *testing.T, val, expected interface{}) {
	t.Helper()

	if val != expected {
		t.Errorf("value (%+v) was not like expected (%+v)", val, expected)
	}
}

func TestCrop(t *testing.T) {
	_, err := Crop(testImage, 1005, 0, 1, 1)
	if err != ErrPointOutOfBounds {
		t.Error("should have returned ErrOutOfBounds, but was ", err)
	}

	_, err = Crop(testImage, 0, -1, 1, 1)
	if err != ErrPointOutOfBounds {
		t.Error("should have returned ErrOutOfBounds, but was ", err)
	}

	_, err = Crop(testImage, 0, 0, 0, 1)
	if err != ErrDimensionOutOfBounds {
		t.Error("should have returned ErrDimensionOutOfBounds, but was ", err)
	}

	_, err = Crop(testImage, 0, 100, 1, 950)
	if err != ErrDimensionOutOfBounds {
		t.Error("should have returned ErrDimensionOutOfBounds, but was ", err)
	}

	res, err := Crop(testImage, 50, 50, 200, 250)
	if err != nil {
		t.Error(err)
	}
	bounds := res.Bounds()
	assert(t, bounds.Dx(), 200)
	assert(t, bounds.Dy(), 250)
}

func TestSlice(t *testing.T) {
	const slices = 5

	res, err := Slice(testImage, slices)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 5 {
		t.Errorf("array size was %d instead of 5", len(res))
	}

	expectedSizeX := testRect.Dx() / slices
	expectedSizeY := testRect.Dy() / slices

	for _, row := range res {
		if len(row) != 5 {
			t.Errorf("row size was %d instead of 5", len(row))
		}

		for _, slice := range row {
			bounds := slice.Bounds()
			assert(t, bounds.Dx(), expectedSizeX)
			assert(t, bounds.Dy(), expectedSizeY)
		}
	}
}
