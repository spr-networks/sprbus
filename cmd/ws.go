package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

func ConnectWebsocket(addr string, authString string, callback func(string, string)) {
	//todo move to main
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws_events_all"}
	//log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	authenticated := false

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			if !authenticated && string(message) == "success" {
				authenticated = true
				continue
			}

			if !authenticated {
				log.Panic("not authenticated")
			}

			if len(message) == 0 {
				return
			}

			topic := gjson.Get(string(message), "Type").String()
			value := gjson.Get(string(message), "Data").String()
			callback(topic, value)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	c.WriteMessage(websocket.TextMessage, []byte(authString))

	for {
		select {
		case <-done:
			return
			/*
				case t := <-ticker.C:
					err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
					if err != nil {
						log.Println("write:", err)
						return
					}
			*/
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				//log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
