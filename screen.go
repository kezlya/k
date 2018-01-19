package k

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"log"
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
	if s.layers[2] == nil {
		s.layers[2] = l
		return
	} else if s.layers[1] == nil {
		s.layers[1] = s.layers[2]
		s.layers[2] = l
		return
	} else if s.layers[0] != nil {
		s.layers[0].removed = true
	}
	s.layers[0] = s.layers[1]
	s.layers[1] = s.layers[2]
	s.layers[2] = l
	return
}

func (s *Screen) Remove(l *Layer) {
	log.Println("Removing layer from the screen")
	l.removed = true
	for i, _l := range s.layers {
		if _l == l {
			s.layers[i] = nil
			return
		}
	}
}
func (s *Screen) RemoveAll() {
	log.Println("Removing All layers from the screen")
	for _, _l := range s.layers {
		if _l != nil {
			s.Remove(_l)
		}
	}
}

func (s *Screen) Display(width, height int) *image.RGBA {
	o := image.Point{0, 0}
	b := image.Rect(0, 0, width, height)
	d := image.NewRGBA(b)
	draw.Draw(d, b, &image.Uniform{color.White}, o, draw.Src)

	for _, l := range s.layers {
		if l != nil {
			draw.Draw(d, b, l.Still, o, draw.Over)
		}
	}

	if s.grid == FOUR {
		w, h := width/2, height/2
		o2, o3, o4 := image.Pt(0, -h), image.Pt(-w, -h), image.Pt(-w, 0)
		sd := resize.Thumbnail(uint(w), uint(h), d, resize.Bicubic).(*image.RGBA)
		draw.Draw(d, b, sd, o, draw.Over)
		draw.Draw(d, b, sd, o2, draw.Over)
		draw.Draw(d, b, sd, o3, draw.Over)
		draw.Draw(d, b, sd, o4, draw.Over)
	} else if s.grid == EIGHT {
		w, h := width/4, height/4
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				sd := resize.Thumbnail(uint(w), uint(h), d, resize.Bicubic).(*image.RGBA)
				draw.Draw(d, b, sd, image.Pt(-w*i, -h*j), draw.Over)
			}
		}
	}

	return d
}

func (s *Screen) GridTo(size DisplayGrid) {
	s.grid = size
}

func blank() *image.RGBA {
	sq := image.Rectangle{image.Point{0, 0}, image.Point{100, 100}}
	return image.NewRGBA(sq)
}
