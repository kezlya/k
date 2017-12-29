package k

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"
	"net/http"
	"github.com/nfnt/resize"
)

type Layer struct {
	// isLive   bool maybe will use it in the future
	backup *image.RGBA
	Still  *image.RGBA
}

func getRandomLayer(width, height int) *Layer {
	return &Layer{false, nil, randomPixels(width, height)}
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

func layerFromImage(keyword string) *k.Layer {
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
							i = 1000 //TODO: refactor
						}
					}
				}
			})
		})
	}

	return &k.Layer{false, nil, convertYCbCr_RGBA(img.(*image.YCbCr))}
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


// Voided function to run whenewer they finish

func scaleDown(layer *k.Layer, rate time.Duration, loop bool) {
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