package randgur

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type FoundCallbackT func(string, []byte)
type Randgur struct {
	HttpClient  http.Client
	workers     chan bool
	Concurrency uint8

	FoundCallbacks []FoundCallbackT
}

func writeToFile(imgID string, body []byte) {
	f, err := os.Create(fmt.Sprintf("./images/%s.jpg", imgID))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(body)
}

var nameChars []byte
var randomSource rand.Source = rand.NewSource(time.Now().UnixNano())
var random = rand.New(randomSource)

func (c *Randgur) Start() {
	c.workers = make(chan bool, c.Concurrency)
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
			time.Sleep(50 * time.Millisecond)
		}
	}
}
func (c *Randgur) Stop() {

}
func (c *Randgur) RegisterCallback(cb FoundCallbackT) {
	if c.FoundCallbacks == nil {
		c.FoundCallbacks = make([]FoundCallbackT, 0)
	}
	c.FoundCallbacks = append(c.FoundCallbacks, cb)
}
func init() {
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
	resp, err := c.HttpClient.Get(fmt.Sprintf("https://i.imgur.com/%s.jpg/", imgID))
	if err != nil || resp.StatusCode == 302 {
		return
	}
	defer resp.Body.Close()

	imageBytes, _ := ioutil.ReadAll(resp.Body)
	for _, cb := range c.FoundCallbacks {
		cb(imgID, imageBytes)
	}
}
