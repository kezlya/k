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
	"strings"
	"time"
)

const displayWidth, displayHeight, quality = 500, 500, 80

var config map[string]string
var words *k.Stack



func main() {
	loadConfig()

	words = k.NewStack()

	screen := k.Screen{}

	go BestEffectSoFar(&screen)

	//go playGroud(&screen)

	//go RecoverDamage(&screen)

	//go analogNumber(&screen)

	//go listingAndShow(&screen)

	//time.Sleep(10000 * time.Millisecond)

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
		kv := strings.Split(s.Text(), "=")
		if len(kv) == 2 {
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
		jpeg.Encode(w, screen.Display(displayWidth, displayHeight), &jpeg.Options{quality})
		sp := r.URL.Query().Get("word")
		if sp != "" && sp != "undefined" {
			words.Push(&k.Node{strings.Replace(strings.TrimSpace(sp)," ","+",100)})
		}
	})
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func playGroud(screen *k.Screen) {
	//screen.GridTo(k.FOUR)
	//layer3 := k.LayerFrom(k.OnlineImage("http://thedailyrecord.com/files/2011/11/orioles-bird.png"))
	for i := 0; i < 10; i++ {
		layer3 := k.LayerFrom(k.GoogleImage("Flowers",-1))
		go layer3.FadeIn(5)
		screen.Add(layer3)
		time.Sleep(2000 * time.Millisecond)

	}
	screen.RemoveAll()
	return
}

func BestEffectSoFar(screen *k.Screen) {
	//screen.GridTo(k.FOUR)
	for i := 0; i < 10; i++ {
		layer3 := k.LayerFrom(k.GoogleImage("png",-1))
		go layer3.FadeIn(5)
		go layer3.ScaleUp(200,500,false)
		screen.Add(layer3)
		time.Sleep(2000 * time.Millisecond)

	}
	screen.RemoveAll()
	return
}

func RecoverDamage(screen *k.Screen) {
	layer3 := k.LayerFrom(k.OnlineImage("http://lsusmath.rickmabry.org/rmabry/knots/newfauxtrefoil2-500x500.jpg"))
	screen.Add(layer3)
	layer1 := k.LayerFrom(k.RandomAlpha(displayWidth,displayHeight))
	screen.Add(layer1)
	go layer1.FadeOut(200)
	return
}

func listingAndShow(screen *k.Screen) {

	//screen.GridTo(k.EIGHT)
	for {
		time.Sleep(400 * time.Millisecond)
		if w := words.Pop(); w != nil {
			l := k.LayerFrom(k.GoogleImage(w.Value, -1))
			go l.ScaleUp(30, 800, true)
			screen.Add(l)
		}

	}

	screen.RemoveAll()
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
