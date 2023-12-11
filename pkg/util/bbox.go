package util

import (
	"errors"
	"image"
	"strconv"
	"strings"
)

var (
	ErrInvalidBBoxParam = errors.New("Invalid bbox param")
)

func ParseBBoxString(bbox string) (*image.Rectangle, error) {
	var ints []int
	for _, part := range strings.Split(bbox, ",") {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, ErrInvalidBBoxParam
		}
		ints = append(ints, num)
	}

	if len(ints) != 4 {
		return nil, ErrInvalidBBoxParam
	}

	x := ints[0]
	y := ints[1]
	w := ints[2]
	h := ints[3]

	rect := image.Rect(x, y, x+w, y+h)

	return &rect, nil
}
