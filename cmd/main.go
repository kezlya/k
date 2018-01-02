package main

import (
	"github.com/kezlya/k"
	"image/jpeg"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const displayWidth, displayHeight, quality = 500, 500, 80

func main() {
	screen := k.Screen{}

	//playGroud(&screen)

	go analogNumber(&screen)

	startServer(&screen)
}

func startServer(screen *k.Screen) {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/stream.jpg", func(w http.ResponseWriter, r *http.Request) {
		jpeg.Encode(w, screen.Display(displayWidth, displayHeight), &jpeg.Options{quality})
	})
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func playGroud(screen *k.Screen) {
	layer1 := k.LayerFrom(k.RandomPixels(500, 500))
	go layer1.ScaleDown(150, true)
	screen.Add(layer1)

	layer3 := k.LayerFrom(k.OnlineImage("http://2fatnerds.com/wp-content/uploads/2013/10/giant-manta-ray.png"))
	go layer3.ScaleDown(33, true)
	screen.Add(layer3)

	rand.Seed(time.Now().UnixNano())
	guess := rand.Intn(100)
	layer2 := k.LayerFrom(k.GoogleImage("Number+"+strconv.Itoa(guess), -1))
	go layer2.ScaleDown(77, true)
	screen.Add(layer2)
}

func analogNumber(screen *k.Screen) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	l := k.LayerFrom(k.GoogleImage("Number+"+strconv.Itoa(n), -1))
	screen.Add(l)
	go l.ScaleDown(77, true)

	for {
		n = rand.Intn(100)
		l.Still = k.GoogleImage("Number+"+strconv.Itoa(n), -1)
		time.Sleep(500 * time.Millisecond)
	}
}
