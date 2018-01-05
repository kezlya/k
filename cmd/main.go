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

	playGroud(&screen)

	//go analogNumber(&screen)

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

	//layer3 := k.LayerFrom(k.RandomPixels(500,500))
	//layer3 := k.LayerFrom(k.OnlineImage("http://thedailyrecord.com/files/2011/11/orioles-bird.png"))
	layer3 := k.LayerFrom(k.GoogleImage("ny+pogadi", 17))
	go layer3.ScaleUp(33, 700,true)
	screen.Add(layer3)
	screen.GridTo(k.FOUR)

}

func analogNumber(screen *k.Screen) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	l := k.LayerFrom(k.GoogleImage("Number+"+strconv.Itoa(n), -1))
	screen.Add(l)
	screen.GridTo(k.FOUR)
	go l.ScaleDown(77, true)

	for {
		n = rand.Intn(100)
		l.Still = k.GoogleImage("Number+"+strconv.Itoa(n), -1)
		time.Sleep(500 * time.Millisecond)
	}
}
