package main

import (
	"net/http"
	"time"
)

func main() {
	time.Sleep(10 * time.Second)
	sendSimulatorDone()
}

func sendSimulatorDone() {
	http.PostForm("http://localhost:8080/simulator/done", map[string][]string{})
}
