package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coder/websocket"
)

func main() {
	if err := doMain(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func doMain() error {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServerFS(os.DirFS(".")))

	mux.HandleFunc("/websocket", handleWebsocket)

	http.ListenAndServe("localhost:9124", mux)

	return nil
}

func handleWebsocket(rw http.ResponseWriter, req *http.Request) {
	conn, err := websocket.Accept(rw, req, nil)
	if err != nil {
		log.Printf("failed to accept: %s", err)
		return
	}
	defer conn.CloseNow()
	readCtx := conn.CloseRead(context.Background())

	log.Printf("connection established with %s", req.RemoteAddr)

	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-readCtx.Done():
			if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
				log.Printf("failed to close: %s", err)
			}
			log.Printf("connection closed")
			return
		case <-t.C:
			if err := conn.Write(context.Background(), websocket.MessageText, []byte("reload")); err != nil {
				log.Printf("failed to write: %s", err)
				return
			}
		}
	}
}
