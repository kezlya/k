package k

import "image"

type Layer struct {
	isLive   bool
	original *image.RGBA
	current  *image.RGBA
}