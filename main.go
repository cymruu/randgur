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
	workers     chan bool
	concurrency uint8
}

var nameChars []byte
var randomSource rand.Source = rand.NewSource(time.Now().UnixNano())
var random = rand.New(randomSource)

func (c *Randgur) Start() {
	c.workers = make(chan bool, c.concurrency)
	for i := 0; i < cap(c.workers); i++ {
		c.workers <- true
	}
	for {
		select {
		case <-c.workers:
			// fmt.Printf("Starting new worker\n")
			go func() {
				defer func() { c.workers <- true }()
				c.GetImage()
			}()
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
func (c *Randgur) GetImage() {
	imgID := c.GuessImageID(7)
	resp, err := c.httpClient.Get(fmt.Sprintf("https://i.imgur.com/%s.jpg/", imgID))
	// fmt.Printf("req for %s => %d", imgID, resp.StatusCode)
	if err != nil || resp.StatusCode == 302 {
		return
	}
	defer resp.Body.Close()
	f, err := os.Create(fmt.Sprintf("./images/%s.jpg", imgID))
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Printf("Found %s \n", imgID)
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
	client.Start()

	fmt.Printf("Statystyki")
}
