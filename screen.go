package k

import (
	"image"
	"image/draw"
	"log"
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
func (s *Screen) Display() *image.RGBA {
	//TODO: merge layers and return result. For now merging only two layers.

	if s.layers[0].current != nil && s.layers[1].current != nil {

		//TODO: test performance and refactor
		sp2 := image.Point{s.layers[0].current.Bounds().Dx(), 0}
		r2 := image.Rectangle{sp2, sp2.Add(s.layers[1].current.Bounds().Size())}
		r := image.Rectangle{image.Point{0, 0}, r2.Max}

		rgba := image.NewRGBA(r)
		draw.Draw(rgba, s.layers[0].current.Bounds(), s.layers[0].current, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, r2, s.layers[1].current, image.Point{0, 0}, draw.Src)

		return rgba
	} else {
		return blank()
	}
}

func blank() *image.RGBA {
	sq := image.Rectangle{image.Point{0, 0}, image.Point{100, 100}}
	return image.NewRGBA(sq)
}
