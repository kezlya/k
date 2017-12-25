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

var img *image.RGBA

type Animation struct {
	Length int
	Img    *image.RGBA
}

func main() {
	sq := image.Rectangle{image.Point{0, 0}, image.Point{500, 500}}
	img = image.NewRGBA(sq)

	go sartAnimation()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/stream.png", func(w http.ResponseWriter, r *http.Request) {
		png.Encode(w, img)
	})
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func trackMouse() {
	//for {
	//	hh := mouse.Event{}
	//	fmt.Println("pos:", hh.X, hh.Y)
	//	time.Sleep(50 * time.Microsecond)
	//}
}

func randomPixels() {
	for i := 0; i < 255; i++ {

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
		}
}

func sartAnimation() {
	//animate forever
	randomPixels()
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


		img = scaleDown(5)
		time.Sleep(50 * time.Nanosecond)

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
func scaleDown(times int) *image.RGBA {
	bb := resize.Thumbnail(uint(100),uint(100),img,resize.Bicubic)
	return bb.(*image.RGBA)
}

func scaleUp(speed, times int) {

}