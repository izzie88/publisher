package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address") // define the address to be used for the URL for the connecting websocket server

type msg string

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	var t msg
	t = "I am sending you this message"
	u := url.URL{ ///define the URL to be used by the websocket defaultdailer dial function
		Scheme: "ws",
		Host:   *addr,
		Path:   "/",
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil) /// Dial into the websocket server
	log.Printf("Dialing %s ...... completed", u.String())
	if err != nil {
		log.Fatal("dial:", err) /// handle error thrown by the websocket dial
	}
	//defer c.Close()
	err = c.WriteMessage(websocket.TextMessage, []byte(t))
	if err != nil {
		log.Println("write:", err)
	}
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
		}
	}
}
