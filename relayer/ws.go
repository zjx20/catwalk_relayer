package relayer

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 4) / 10
)

func reader(ws *websocket.Conn, upConn net.Conn) {
	defer func() {
		ws.Close()
		upConn.Close()
	}()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		typ, buf, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if typ != websocket.BinaryMessage {
			log.Printf("required BinaryMessage, got %v", typ)
			break
		}
		_, err = upConn.Write(buf)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func writer(ws *websocket.Conn, upConn net.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
		upConn.Close()
	}()
	const BufSize = 4096
	bufPool := sync.Pool{
		New: func() interface{} {
			return make([]byte, BufSize)
		},
	}
	dataCh := make(chan []byte, 8)
	go func() {
		for {
			buf := bufPool.Get().([]byte)
			buf = buf[:cap(buf)]
			n, err := upConn.Read(buf)
			if err != nil {
				log.Println(err)
				dataCh <- nil
				return
			}
			dataCh <- buf[:n]
		}
	}()
	for {
		select {
		case buf := <-dataCh:
			if buf == nil {
				return
			}
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.BinaryMessage, buf); err != nil {
				return
			}
			bufPool.Put(buf)
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Server Internal Error", 500)
}

func serveWs(w http.ResponseWriter, r *http.Request, upstream string) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Expect Upgrade", 400)
		return
	}

	d := net.Dialer{
		Timeout: 10 * time.Second,
	}
	upConn, err := d.Dial("tcp", upstream)
	if err != nil {
		http.Error(w, "Internal Error", 501)
		log.Println(err)
		return
	}

	go writer(conn, upConn)
	reader(conn, upConn)
}

func WsRelay(port int, wsPath string, upstream string) {
	log.Printf("Relay Websocket on port %d, upstream %s\n", port, upstream)
	http.HandleFunc("/", serveHome)
	http.HandleFunc(wsPath, func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, upstream)
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
