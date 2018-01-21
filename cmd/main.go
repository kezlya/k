package main

import (
	"bufio"
	"fmt"
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
var words *Stack

type Node struct {
	Value string
}

func (n *Node) String() string {
	return fmt.Sprint(n.Value)
}

// NewStack returns a new stack.
func NewStack() *Stack {
	return &Stack{}
}

// Stack is a basic LIFO stack that resizes as needed.
type Stack struct {
	nodes []*Node
	count int
}

// Push adds a node to the stack.
func (s *Stack) Push(n *Node) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *Stack) Pop() *Node {
	if s.count == 0 {
		return nil
	}
	s.count--
	return s.nodes[s.count]
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*Node, size),
		size:  size,
	}
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	nodes []*Node
	size  int
	head  int
	tail  int
	count int
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Node) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*Node, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

func main() {
	loadConfig()

	words = NewStack()

	screen := k.Screen{}

	go playGroud(&screen)

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
			words.Push(&Node{strings.Replace(strings.TrimSpace(sp)," ","+",100)})
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
	//for i := 0; i < 10; i++ {
		layer3 := k.LayerFrom(k.GoogleImage("Flowers",-1))
		go layer3.FadeIn(10)
		screen.Add(layer3)
		time.Sleep(10000 * time.Millisecond)

	//}
	//screen.RemoveAll()
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
			log.Println(words.count,w.Value)
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
