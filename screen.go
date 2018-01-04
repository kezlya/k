package k

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"github.com/nfnt/resize"
	"fmt"
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
		gridBG := image.NewRGBA(image.Rect(0, 0, width, height))
//fmt.Println(gridBG.Bounds())

		sw := width/2
		sh := height/2
		//rc := image.Rect(0,0,sw,sh)
		pt1, pt2, pt3, pt4 := image.Pt(0,0), image.Pt(0, -sh), image.Pt(-sw,-sh), image.Pt(-sw,0)
		sRgba := resize.Thumbnail(uint(sw), uint(sh), rgba, resize.Bicubic).(*image.RGBA)

		fmt.Println(rgba.Bounds())


		draw.Draw(gridBG,gridBG.Bounds(),sRgba,pt1, draw.Over)
		draw.Draw(gridBG,gridBG.Bounds(),sRgba,pt2, draw.Over)
		draw.Draw(gridBG,gridBG.Bounds(),sRgba,pt3, draw.Over)
		draw.Draw(gridBG,gridBG.Bounds(),sRgba,pt4, draw.Over)



		b := image.Rect(0, 0, width, height)
		fmt.Println(image.Pt(0, 0).Sub(sRgba.Bounds().Max))
		//b.in
		//p := image.Pt(0, 30)
		// Note that even though the second argument is b,
		// the effective rectangle is smaller due to clipping.
		draw.Draw(rgba, b, sRgba, image.Pt(0, 0).Sub(sRgba.Bounds().Max), draw.Over)
		//dirtyRect := b.Intersect(image.Rect(b.Min.X, b.Max.Y-20, b.Max.X, b.Max.Y))
		return gridBG





		//return gridBG
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
