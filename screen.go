package k

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"github.com/nfnt/resize"
)

type DisplayGrid int

const (
	ONE DisplayGrid = 1 << iota
	TWO
	FOUR
	EIGHT
	SIXTEEN
)

type Screen struct {
	layers [3]*Layer
	grid   DisplayGrid
}

func (s *Screen) Add(l *Layer) {
	log.Println("adding layer to the screen")

	//TODO: implement queue or play with recursive
	if s.layers[0] == nil {
		s.layers[0] = l
		return
	} else if s.layers[1] == nil {
		s.layers[1] = s.layers[0]
		s.layers[0] = l
	} else {
		s.layers[2] = s.layers[1]
		s.layers[1] = s.layers[0]
		s.layers[0] = l
	}
}

func (s *Screen) Remove(l *Layer) {
	log.Println("Removing layer from the screen")

	for i, _l := range s.layers {
		if _l == l {
			s.layers[i] = nil
			return
		}
	}
}

func (s *Screen) Display(width, height int) *image.RGBA {
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)

	if s.layers[0] != nil {
		draw.Draw(rgba, rgba.Bounds(), s.layers[0].Still, s.layers[0].Still.Bounds().Min, draw.Over)
	}

	if s.layers[1] != nil {
		draw.Draw(rgba, rgba.Bounds(), s.layers[1].Still, s.layers[1].Still.Bounds().Min, draw.Over)
	}

	if s.layers[2] != nil {
		draw.Draw(rgba, rgba.Bounds(), s.layers[2].Still, s.layers[2].Still.Bounds().Min, draw.Over)
	}

	if s.grid == FOUR{
		sRgba := resize.Thumbnail(uint(width/2), uint(height/2), rgba, resize.Bicubic).(*image.RGBA)
		rgba = image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(rgba, rgba.Bounds(), sRgba, image.Point{0,0}, draw.Over)
		draw.Draw(rgba, rgba.Bounds(), sRgba, image.Point{width/2,0}, draw.Over)
		draw.Draw(rgba, rgba.Bounds(), sRgba, image.Point{0,height}, draw.Over)
		draw.Draw(rgba, rgba.Bounds(), sRgba, image.Point{width/2,height}, draw.Over)
	}

	return rgba
}

func (s *Screen) GridTo(size DisplayGrid) {
	s.grid = size
}

func blank() *image.RGBA {
	sq := image.Rectangle{image.Point{0, 0}, image.Point{100, 100}}
	return image.NewRGBA(sq)
}
