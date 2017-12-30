package k

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Layer struct {
	// isLive   bool maybe will use it in the future
	backup *image.RGBA
	Still  *image.RGBA
}

func RandomPixels(width, height int) *Layer {
	sq := image.Rectangle{
		image.Point{0, 0},
		image.Point{width, height}}
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
	return &Layer{Still: img}
}

func GoogleImage(keyword string, order int) *Layer {
	var img *image.RGBA

	if order < 1 {
		rand.Seed(time.Now().UnixNano())
		order = rand.Intn(20)
	}

	doc, err := goquery.NewDocument("https://www.google.com/search?q=" + keyword + "&tbm=isch")
	if err != nil {
		fmt.Printf(err.Error())
		log.Fatal(err)
	}

	doc.Find(".images_table").Each(func(index int, item *goquery.Selection) {
		item.Find("img").Each(func(index2 int, item2 *goquery.Selection) {
			if index2 == order {
				if src, e := item2.Attr("src"); e == true {
					img = loadFromUrl(src)
				}
			}
		})
	})

	if img == nil {
		img = blank()
	}

	return &Layer{Still: img}
}

func OnlineImage(url string) *Layer {
	var img *image.RGBA

	img = loadFromUrl(url)
	if img == nil {
		img = blank()
	}

	return &Layer{Still: img}
}

func convertYCbCr_RGBA(img *image.YCbCr) *image.RGBA {
	b := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
	return m
}


func loadFromUrl(url string) *image.RGBA {
	fmt.Println("load: ", url)
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		log.Println(res.StatusCode , " status code from the url ", url)
	}
	defer res.Body.Close()
	m, _, _ := image.Decode(res.Body)

	_, ok := m.(*image.YCbCr)
	if !ok{
		return convertYCbCr_RGBA(m.(*image.YCbCr))
	}

	_, ok = m.(*image.RGBA)
	if !ok{
		return nil
	}

	return m.(*image.RGBA)
}

func scaleUp(speed, times int) {

}

func (s *Layer) ScaleDown(rate time.Duration, loop bool) {
	if loop {
		s.backup = s.Still
	}
	for {
		time.Sleep(rate * time.Millisecond)
		size := s.Still.Rect.Size()
		if size.X > 1 && size.Y > 1 {
			bb := resize.Thumbnail(uint(size.X-5), uint(size.Y-5), s.Still, resize.Bicubic)
			s.Still = bb.(*image.RGBA)
			time.Sleep(rate * time.Millisecond)
		} else {
			break
		}
	}
	if loop {
		s.Still = s.backup
		s.ScaleDown(rate, loop)
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
