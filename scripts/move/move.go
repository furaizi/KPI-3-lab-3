package main

import (
	"bytes"
	"net/http"
	"time"
	"fmt"
)

func main() {
	data := []byte(
		`white
bgrect 0.15 0.15 0.35 0.35
green
figure 0.25 0.25
update`)
	_, err := http.Post("http://localhost:17000/", "application/text", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	moveData := []byte(`move 0.01 0.01
update`)
	for {
		resp, err := http.Post("http://localhost:17000/", "application/text", bytes.NewBuffer(moveData))
		if err != nil {
			fmt.Println("Error occurred:", err)
            time.Sleep(100 * time.Millisecond)
			continue
		}
		resp.Body.Close()
        time.Sleep(100 * time.Millisecond)
	}
}