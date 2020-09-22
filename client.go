package gowsclient

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10
)

type Message struct {
	Mt      int
	Message []byte
}

type client struct {
	conn *websocket.Conn

	wsUrl  *url.URL
	header *http.Header

	sendCh chan *Message

	OnConnect func()
	OnReceive func(msg *Message)
}

func New(rawUrl string) (*client, error) {

	return NewWithHeader(rawUrl, nil)
}

func NewWithHeader(rawUrl string, header *http.Header) (*client, error) {

	wsUrl, err := url.Parse(rawUrl)

	if err != nil {
		return nil, err
	}

	c := &client{
		header: header,
		wsUrl:  wsUrl,
		sendCh: make(chan *Message),
	}

	return c, nil
}

func (c *client) reconnect() {

	c.Close()
	c.connect()
}

func (c *client) Close() {

	if c.conn != nil {
		c.conn.Close()
	}

}
func (c *client) connect() {

	var header http.Header

	if c.header != nil {
		header = *c.header
	}
	ws, _, err := websocket.DefaultDialer.Dial(c.wsUrl.String(), header)

	c.conn = ws

	if err == nil {
		log.Printf("Dial: connection was successfully established with %s\n", c.wsUrl.String())

		c.start()

		if c.OnConnect != nil {
			c.OnConnect()
		}
		return
	}

	log.Println("Could not connect:", err.Error())

	after := time.After(time.Second * 2)

	<-after

	c.connect()

}

func (c *client) WriteMessage(msg *Message) {

	c.sendCh <- msg

}

func (c *client) handleWrite() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()

	}()

	for {

		select {
		case msg := <-c.sendCh:


			if c.conn == nil {
				return
			}

			err := c.conn.WriteMessage(msg.Mt, msg.Message)

			if err != nil {

				return

			}
		case <-ticker.C:


			if c.conn == nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}

	}
}

func (c *client) start() {

	go c.loopRead()
	go c.handleWrite()

}

func (c *client) loopRead() {
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	defer c.reconnect()

	for {
		mt, message, err := c.conn.ReadMessage()

		if err != nil {

			log.Println("Read error :", err.Error())
			break
		}

		if c.OnReceive != nil {
			c.OnReceive(&Message{
				Mt:      mt,
				Message: message,
			})
		}

	}
}
