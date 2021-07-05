package relayer

import (
	"fmt"
	"log"
	"net"
	"time"
)

func pipe(c1 net.Conn, c2 net.Conn) {
	defer func() {
		c1.Close()
		c2.Close()
	}()
	c1.(*net.TCPConn).ReadFrom(c2)
}

func TcpRelay(port int, upstream string) {
	log.Printf("Relay Tcp on port %d, upstream %s\n", port, upstream)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		d := net.Dialer{
			Timeout: 10 * time.Second,
		}
		upConn, err := d.Dial("tcp", upstream)
		if err != nil {
			log.Println(err)
			continue
		}
		go pipe(conn, upConn)
		go pipe(upConn, conn)
	}
}
