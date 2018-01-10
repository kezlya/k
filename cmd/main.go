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
var listening []string

func main() {
	listening = append(listening,"")
	loadConfig()

	screen := k.Screen{}

	//playGroud(&screen)

	go listingAndShow(&screen)

	//go analogNumber(&screen)

	//startListining()

	startServer(&screen)
}

func loadConfig() {
	config = make(map[string]string)
	f, err := os.Open("config.ignore")
	if err != nil {
		log.Fatalln("Can't find config.ignore file: ", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		kv := strings.Split(s.Text(),"=")
		if len(kv) == 2{
			config[kv[0]] = kv[1]
		}
	}
	if s.Err() != nil {
		log.Fatalln("Problems with some variables in config.ignore: ", s.Err())
	}
}

func startServer(screen *k.Screen) {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/stream.jpg", func(w http.ResponseWriter, r *http.Request) {
		sp := r.URL.Query().Get("speech")
		if sp !="" && sp != listening[len(listening)-1]{
			listening = append(listening,sp)
		}
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

func listingAndShow(screen *k.Screen){
	var last string
	l := k.LayerFrom(k.GoogleImage("talk",-1))
	screen.Add(l)
	screen.GridTo(k.EIGHT)
	go l.ScaleUp(30, 800, true)

	for {
		last = strings.Replace(listening[len(listening)-1]," ","+",-1)
		l.Still = k.GoogleImage(last, -1)
		time.Sleep(2 * time.Second)
		log.Println(len(listening),last)
	}
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
