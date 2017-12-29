package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kezlya/k"
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
	"github.com/nfnt/resize"
)

func main() {
	screen := k.Screen{}

	layer1 := getRandomLayer(100, 100)
	//go scaleDown(layer1, 150, true)
	screen.Add(layer1)

	rand.Seed(time.Now().UnixNano())
	guess := rand.Intn(100)
	layer2 := layerFromImage("Number+" + strconv.Itoa(guess))
	//go scaleDown(layer2, 77, true)
	screen.Add(layer2)

	//layer3 := getRandomLayer(300, 300)
	//screen.Add(layer3.current)

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

