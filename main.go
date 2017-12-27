package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Layer struct {
	original *image.RGBA
	current  *image.RGBA
}

type Screen struct {
	layers [3]*Layer
}

func (s Screen) AddLayer(l *Layer) {
	fmt.Println("adding layer to the screen")

	//TODO: implement queue
	if s.layers[0] == nil {
		fmt.Println("1")
		s.layers[0] = l
		return
	} else if s.layers[1] == nil {
		fmt.Println("2")

		s.layers[1] = s.layers[0]
		s.layers[0] = l
	} else {
		fmt.Println("3")

		s.layers[2] = s.layers[1]
		s.layers[1] = s.layers[0]
		s.layers[0] = l
	}
}

func (s Screen) RemoveLayer(l *Layer) {
	fmt.Println("adding layer to the screen")

	for i, _l := range s.layers {
		if _l == l {
			s.layers[i] = nil
			return
		}
	}
}

func (s Screen) Display() *image.RGBA {
	//TODO: merge layers and return result. For now merging only two layers.

	if s.layers[0] != nil && s.layers[1] != nil {

		//TODO: test performance and refactor
		sp2 := image.Point{s.layers[0].current.Bounds().Dx(), 0}
		r2 := image.Rectangle{sp2, sp2.Add(s.layers[1].current.Bounds().Size())}
		r := image.Rectangle{image.Point{0, 0}, r2.Max}
		rgba := image.NewRGBA(r)
		draw.Draw(rgba, s.layers[1].current.Bounds(), s.layers[1].current, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, r2, s.layers[2].current, image.Point{0, 0}, draw.Src)

		return rgba
	} else {
		return randomPixels(100, 100)

	}
}

func main() {
	screen := Screen{[3]*Layer{nil,nil,nil}}

	rand.Seed(time.Now().UnixNano())
	guess := rand.Intn(100)

	layer1 := getRandomLayer()
	screen.AddLayer(layer1)

	layer2 := layerFromImage("Number+" + strconv.Itoa(guess))
	go sartAnimation(layer2)
	screen.AddLayer(layer2)

	fmt.Printf("", screen)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/stream.jpg", func(w http.ResponseWriter, r *http.Request) {
		jpeg.Encode(w, screen.Display(), &jpeg.Options{80})
	})
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getRandomLayer() *Layer {
	return &Layer{current: randomPixels(500, 500)}
}

func randomPixels(x, y int) *image.RGBA {
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

func layerFromImage(keyword string) *Layer {
	var img image.Image

	doc, err := goquery.NewDocument("https://www.google.com/search?q=" + keyword + "&tbm=isch")
	if err != nil {
		fmt.Printf(err.Error())
		log.Fatal(err)
	}

	// slow on purpose it influence analog like
	// try 3 times

	for i := 1; i <= 3; i++ {
		rand.Seed(time.Now().UnixNano())
		guess := rand.Intn(10)
		fmt.Println(guess)
		doc.Find(".images_table").Each(func(index int, item *goquery.Selection) {
			item.Find("img").Each(func(index2 int, item2 *goquery.Selection) {
				if index2 == guess {
					if src, e := item2.Attr("src"); e == true {
						img = loadJpegFromUrl(src)
						if img != nil {
							i = 1000
						}
					}
				}
			})
		})
	}

	return &Layer{current: convertYCbCr_RGBA(img.(*image.YCbCr))}
}

func convertYCbCr_RGBA(img *image.YCbCr) *image.RGBA {
	b := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
	return m
}

func loadJpegFromUrl(url string) image.Image {
	fmt.Println("load: ", url)
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		// handle errors
	}
	defer res.Body.Close()
	m, _, _ := image.Decode(res.Body)
	return m
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

		scaleDown(layer, 77, true)

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
	if loop {
		layer.original = layer.current
	}
	for {
		time.Sleep(rate * time.Millisecond)
		size := layer.current.Rect.Size()
		if size.X > 1 && size.Y > 1 {
			bb := resize.Thumbnail(uint(size.X-5), uint(size.Y-5), layer.current, resize.Bicubic)
			layer.current = bb.(*image.RGBA)
			time.Sleep(rate * time.Millisecond)

		} else {
			break
		}
	}
	if loop {
		layer.current = layer.original
		scaleDown(layer, rate, loop)
	}
}

func scaleUp(speed, times int) {

}
