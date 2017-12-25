package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"time"
	"image/color"
	"github.com/nfnt/resize"
)



type Layer struct {
	original *image.RGBA
	current *image.RGBA
}

//TODO: not sure about it
type Animation struct {
	Length int
	Img    *image.RGBA
}

func main() {
	layer1 := getRandomLayer()


	go sartAnimation(layer1)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/stream.png", func(w http.ResponseWriter, r *http.Request) {
		png.Encode(w, layer1.current)
	})
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getRandomLayer() *Layer {
	return &Layer{current:randomPixels(500,500)}
}

func trackMouse() {
	//for {
	//	hh := mouse.Event{}
	//	fmt.Println("pos:", hh.X, hh.Y)
	//	time.Sleep(50 * time.Microsecond)
	//}
}

func randomPixels(x,y int ) *image.RGBA {
	sq := image.Rectangle{image.Point{0, 0}, image.Point{x, y}}
	var img *image.RGBA
	img = image.NewRGBA(sq)

	for x := 0; x < 500; x++ {
		for y := 0; y < 500; y++ {
			r := uint8(x)
			g := uint8(y)
			b := uint8(rand.Intn(255))
			a := uint8(rand.Intn(255))
			c := color.RGBA{r, g, b, a}
			img.Set(x, y, c)
		}
	}
	return img
}

func sartAnimation(layer *Layer) {
	//animate forever

	for {
		// start adding effects

		//TODO: Convert this colorshifting to functions
		/*for i := 0; i < 255; i++ {

			for x := 0; x < 500; x++ {
				for y := 0; y < 500; y++ {
					r := uint8(x - y - i/2)
					g := uint8(y + x/2 + i)
					b := uint8(rand.Intn(255))
					a := uint8(i)
					c := color.RGBA{r, g, b, a}
					img.Set(x, y, c)
				}
				time.Sleep(50 * time.Nanosecond)
			}
			t := time.Now()
			fmt.Println(t.Second(), t.Nanosecond()) // prints 1000
		}*/

		scaleDown(layer,77, true)


		t := time.Now()
		fmt.Println(t.Second(), t.Nanosecond()) // prints 1000
	}
}

func mirror(n int) int {
	if n > 127 {
		return n - 127
	}
	return n
}

func fade(n, m int) int {
	if n > m {
		return n - m
	}
	return 0
}

func frameNumber(n int) int {
	if n > 4 {
		return 0
	}
	return n + 1
}

// Voided function to run whenewer they finish
func scaleDown(layer *Layer, rate time.Duration, loop bool) {

	if loop{
		layer.original = layer.current
	}

	for {
		time.Sleep(rate * time.Millisecond)
		size := layer.current.Rect.Size()
		if size.X > 1 && size.Y > 1 {
			bb := resize.Thumbnail(uint(size.X-10), uint(size.Y-10), layer.current, resize.Bicubic)
			layer.current = bb.(*image.RGBA)
			time.Sleep(rate * time.Millisecond)

		} else {
			break
		}
	}
	if loop{
		layer.current = layer.original
		scaleDown(layer,rate,loop)
	}
}

func scaleUp(speed, times int) {

}
