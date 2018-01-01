package k

import (
	"image"
	"image/draw"
	"log"
	"image/color"
)

type Screen struct {
	layers [3]*Layer
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

//TODO: need to add parameters for dimensions
func (s *Screen) Display(with, hight int) *image.RGBA {
	//TODO: merge layers and return result. For now merging only two layers.

	rgba := image.NewRGBA(image.Rect(0,0,with,hight))
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)

	if s.layers[0]!= nil {
		draw.Draw(rgba, rgba.Bounds(), s.layers[0].Still, s.layers[0].Still.Bounds().Min, draw.Over)
	}

	if s.layers[1] != nil {
		draw.Draw(rgba, rgba.Bounds(), s.layers[1].Still, s.layers[1].Still.Bounds().Min, draw.Over)
	}

	if s.layers[2] != nil {
		draw.Draw(rgba, rgba.Bounds(), s.layers[2].Still, s.layers[2].Still.Bounds().Min, draw.Over)
	}

	return rgba
}

func blank() *image.RGBA {
	sq := image.Rectangle{image.Point{0, 0}, image.Point{100, 100}}
	return image.NewRGBA(sq)
}
