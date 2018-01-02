package main

import (
	"github.com/kezlya/k"
	"image/jpeg"
	"log"
	"math/rand"
	"net/http"
	"time"
	"strconv"
)

func main() {
	screen := k.Screen{}

	layer1 := k.RandomPixels(500, 500 )
	go layer1.ScaleDown(150, true)
	screen.Add(layer1)


	layer3 := k.OnlineImage("http://2fatnerds.com/wp-content/uploads/2013/10/giant-manta-ray.png")
	go layer3.ScaleDown( 33, true)
	screen.Add(layer3)

	rand.Seed(time.Now().UnixNano())
	guess := rand.Intn(100)
	layer2 := k.GoogleImage("Number+" + strconv.Itoa(guess),-1)
	go layer2.ScaleDown( 77, true)
	screen.Add(layer2)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/stream.jpg", func(w http.ResponseWriter, r *http.Request) {
		jpeg.Encode(w, screen.Display(500,500), &jpeg.Options{80})
	})
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}