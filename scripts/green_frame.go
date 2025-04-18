package main

import (
	"bytes"
	"net/http"
)

func main() {
	data := []byte(
		`white
bgrect 0.15 0.15 0.35 0.35
green
figure 0.25 0.25
update`)
	resp, err := http.Post("http://localhost:17000/", "application/text", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
