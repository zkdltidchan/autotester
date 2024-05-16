package main

import (
	"flag"
	"net/http"
	"time"
)

func main() {
	// if port is not provided, exit
	port := flag.String("port", "", "port")
	flag.Parse()
	if *port == "" {
		return
	}
	sendInitBroker(*port)
	time.Sleep(15 * time.Second)
	sendDone(*port)
}
func sendInitBroker(port string) {
	http.PostForm("http://localhost:8080/init_broker", map[string][]string{"port": {port}})
}

func sendDone(port string) {
	http.PostForm("http://localhost:8080/done", map[string][]string{"port": {port}})
}
