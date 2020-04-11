package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cymruu/randgur"
)

func writeToFile(imgID string, body []byte) {
	f, err := os.Create(fmt.Sprintf("./images/%s.jpg", imgID))
	fmt.Printf("found %s\n", imgID)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(body)
}
func main() {
	client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}, Timeout: time.Second * 5}
	instance := randgur.Randgur{Concurrency: 5, HttpClient: client}
	instance.RegisterCallback(writeToFile)
	instance.Start()
}
