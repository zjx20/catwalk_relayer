package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/zjx20/catwalk_relayer/relayer"
)

var (
	upstream = flag.String("upstream", "", "e.g. www.google.com:80")
	port     = flag.Int("port", defaultPort(), "listen port")
	ws       = flag.Bool("ws", false, "websocket")
	wsPath   = flag.String("wsPath", "/chat", "ws path")
)

func defaultPort() int {
	port := os.Getenv("PORT")
	if port != "" {
		num, err := strconv.Atoi(port)
		if err != nil {
			log.Fatalf("parse PORT env (%s) failed, err: %v", port, err)
		}
		return num
	}
	return 8080
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if *upstream == "" {
		log.Printf("Error: upstream should not be empty\n")
		return
	}
	if *ws {
		relayer.WsRelay(*port, *wsPath, *upstream)
	} else {
		relayer.TcpRelay(*port, *upstream)
	}
}
