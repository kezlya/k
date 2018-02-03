package main

import (
	"bufio"
	"github.com/kezlya/k"
	"image"
	"image/jpeg"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const displayWidth, displayHeight, quality = 200, 100, 80

var config map[string]string
var words *k.Stack
var display *image.RGBA
var lastResponse int64
var screen *k.Screen

func main() {
	//loadConfig()

	words = k.NewStack()

	screen = &k.Screen{}

	//go BestEffectSoFar(&screen)

	go playGroud()
	go screen.RandomGrid(7)

	//go RecoverDamage(&screen)

	//go analogNumber(&screen)

	//go listingAndShow(&screen)

	//time.Sleep(10000 * time.Millisecond)

	startServer()

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

func startServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/stream.jpg", imageStream)
	http.HandleFunc("/number.jpg", imageStream)

	http.HandleFunc("/stream", page("pages/stream.html"))
	http.HandleFunc("/number", page("pages/number.html"))

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func page(n string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, n)
	}
}

func imageStream(w http.ResponseWriter, r *http.Request) {
	lw := "undefined"
	if lastResponse < time.Now().Add(-70*time.Millisecond).UnixNano() {
		display = screen.Display(displayWidth, displayHeight)
		lastResponse = time.Now().UnixNano()
	}
	jpeg.Encode(w, display, &jpeg.Options{quality})
	sp := r.URL.Query().Get("word")
	if sp != "" && sp != "undefined" && sp != lw {
		lw = sp
		words.Push(&k.Node{strings.Replace(strings.TrimSpace(sp), " ", "+", 100), true})
		log.Println("Adding word:", sp)
	}
}


func playGroud() {
	//screen.GridTo(k.FOUR)
	//layer3 := k.LayerFrom(k.OnlineImage("http://thedailyrecord.com/files/2011/11/orioles-bird.png"))
	tries := 1
	var layer3 *k.Layer
	for {
		if w := words.Pop(); w != nil {
			if w.IsVoice {
				tries = 3
			} else {
				tries = 1
			}
			for i := 0; i < tries; i++ {

				if w.IsVoice {
					layer3 = k.LayerFrom(k.GoogleImage(w.Value, -1))

				} else {
					layer3 = k.LayerFrom(k.FlickerImage(w.Value, -1))
				}
				go layer3.FadeIn(20)
				go layer3.RandomEffect()

				screen.Add(layer3)
				time.Sleep(2000 * time.Millisecond)
			}
		} else {
			words.Push(&k.Node{"sky", false})
			words.Push(&k.Node{"mountains", false})
			words.Push(&k.Node{"eye", false})
			words.Push(&k.Node{"space", false})
			words.Push(&k.Node{"lips", false})
			words.Push(&k.Node{"sunshine", false})
			words.Push(&k.Node{"eyes", false})
			words.Push(&k.Node{"Road", false})
			words.Push(&k.Node{"Love", false})
			words.Push(&k.Node{"Kiss", false})
			words.Push(&k.Node{"Highway", false})
			words.Push(&k.Node{"Hand", false})
			words.Push(&k.Node{"Birds", false})
			words.Push(&k.Node{"Berry", false})
			words.Push(&k.Node{"Touch", false})
			words.Push(&k.Node{"Love", false})
			words.Push(&k.Node{"Feels", false})
			words.Push(&k.Node{"More", false})
			words.Push(&k.Node{"Beautiful", false})
		}
	}
	screen.RemoveAll()
	return
}

func BestEffectSoFar() {
	screen.GridTo(k.FOUR)
	for i := 0; i < 10; i++ {
		layer3 := k.LayerFrom(k.GoogleImage("flowers", -1))
		go layer3.BurnOut(7)
		go layer3.ScaleUp(30, 800, false)
		screen.Add(layer3)
		time.Sleep(2000 * time.Millisecond)

	}
	screen.RemoveAll()
	return
}

func RecoverDamage(screen *k.Screen) {
	layer3 := k.LayerFrom(k.OnlineImage("http://lsusmath.rickmabry.org/rmabry/knots/newfauxtrefoil2-500x500.jpg"))
	screen.Add(layer3)
	layer1 := k.LayerFrom(k.RandomAlpha(displayWidth, displayHeight))
	screen.Add(layer1)
	go layer1.FadeOut(200)
	return
}

func listingAndShow(screen *k.Screen) {

	//screen.GridTo(k.EIGHT)
	for {
		time.Sleep(2000 * time.Millisecond)
		if w := words.Pop(); w != nil {
			l := k.LayerFrom(k.GoogleImage(w.Value, -1))
			go l.RandomEffect()
			go l.RandomEffect()
			go l.RandomEffect()
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
