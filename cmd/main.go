package main

import (
	"fmt"
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

	layer1 := k.RandomPixels(100, 100)
	go layer1.ScaleDown(150, true)
	screen.Add(layer1)

	fmt.Println(layer1)

	rand.Seed(time.Now().UnixNano())
	guess := rand.Intn(100)
	layer2 := k.GoogleImage("Number+" + strconv.Itoa(guess),-1)
	go layer2.ScaleDown( 77, true)
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
