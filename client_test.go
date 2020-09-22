package gowsclient

import (
	"log"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNew(t *testing.T) {

	client, err := New("wss://echo.websocket.org")

	if err != nil {
		return
	}

	client.OnConnect = func() {

		client.WriteMessage(&Message{
			Mt:      websocket.TextMessage,
			Message: []byte("Hello!!!"),
		})
	}

	client.OnReceive = func(msg *Message) {

		log.Println(msg.Mt, string(msg.Message))

	}

	go client.connect()

	time.Sleep(time.Second * 20)
}
