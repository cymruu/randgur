package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Randgur struct {
	httpClient  http.Client
	queue       chan string
	images      chan io.ReadCloser
	workers     chan bool
	concurrency uint8
}

var nameChars []byte
var randomSource rand.Source = rand.NewSource(time.Now().UnixNano())
var random = rand.New(randomSource)

func (c *Randgur) Start(guesses uint32) {
	c.queue = make(chan string)
	c.workers = make(chan bool, c.concurrency)
	var i uint32
	for ; i < guesses; i++ {
		go func() {
			c.queue <- c.GuessImageID(7)
		}()
	}
	for i := 0; i < cap(c.workers); i++ {
		c.workers <- true
	}
	for {
		select {
		case imageID := <-c.queue:
			<-c.workers
			fmt.Printf("Got worker for fetching %s.jpg\n", imageID)
			go func() { c.queue <- c.GuessImageID(5) }()
			go func(imageID string) {
				defer func() { c.workers <- true }()
				c.GetImage(imageID)
			}(imageID)
		default:
			fmt.Printf(".")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
func (c *Randgur) Stop() {

}
func Initialize() {
	// uppercase A-Z
	for i := 65; i <= 90; i++ {
		nameChars = append(nameChars, byte(i))
	}
	//lowercase a-z
	for i := 97; i <= 122; i++ {
		nameChars = append(nameChars, byte(i))
	}
	//0-9
	for i := 48; i < 57; i++ {
		nameChars = append(nameChars, byte(i))
	}
}
func (c *Randgur) GuessImageID(size uint8) string {
	guess := make([]byte, size)
	for i := uint8(0); i < size; i++ {
		guess[i] = nameChars[random.Intn(len(nameChars))]
	}
	return string(guess)
}
func (c *Randgur) GetImage(imgID string) {
	resp, err := c.httpClient.Get(fmt.Sprintf("https://i.imgur.com/%s.jpg/", imgID))
	fmt.Printf("req for %s => %d", imgID, resp.StatusCode)
	if err != nil || resp.StatusCode == 302 {
		return
	}
	defer resp.Body.Close()
	f, err := os.Create(fmt.Sprintf("./images/%s.jpg", imgID))
	if err != nil {
		return
	}
	defer f.Close()
	io.Copy(f, resp.Body)
}
func main() {

	Initialize()
	client := Randgur{
		httpClient: http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}},
		concurrency: 5}
	client.Start(5)
}
