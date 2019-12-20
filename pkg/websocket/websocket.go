package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.ximalaya.com/x-fm/tail/pkg/types"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

)

var (
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { _ = ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, fileName string) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	config := types.Config{
		Follow:      true,
		MaxLineSize: 100,
		ReOpen:true,
		//Pipe:true,
	}
	tf, err := types.TailFile(fileName, config)
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case line := <-tf.Lines:
			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.TextMessage, []byte(line.Text)); err != nil {
				return
			}
		case <-pingTicker.C:
			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	//log file name
	fileName:=r.FormValue("filePath")

	go writer(ws, fileName)
	reader(ws)
}

func Start(addr string) {
	http.HandleFunc("/ws", serveWs)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
