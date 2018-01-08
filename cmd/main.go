package main

import (
	"bufio"
	"github.com/kezlya/k"
	"image/jpeg"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"
)

const displayWidth, displayHeight, quality = 500, 500, 80

var config map[string]string

func main() {
	screen := k.Screen{}

	//playGroud(&screen)

	go analogNumber(&screen)

	startServer(&screen)
}

func loadConfig() {
	f, err := os.Open("config.ignore")
	if err != nil {
		log.Fatalln("Can't find config.ignore file: ", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		_kv := strings.Split(s.Text(),"=")
		if len(_kv) == 2{
			config[_kv[0]] = _kv[1]
		}
	}
	log.Println(s.Err())
}

func startListining() {

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
	go layer3.ScaleUp(33, 700, true)
	screen.Add(layer3)
	screen.GridTo(k.FOUR)

}

func analogNumber(screen *k.Screen) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	l := k.LayerFrom(k.GoogleImage("Number+"+strconv.Itoa(n), -1))
	screen.Add(l)
	screen.GridTo(k.EIGHT)
	go l.ScaleUp(30, 800, true)

	for {
		n = rand.Intn(100)
		l.Still = k.GoogleImage("Number+"+strconv.Itoa(n), -1)
		time.Sleep(2000 * time.Millisecond)
	}
}
