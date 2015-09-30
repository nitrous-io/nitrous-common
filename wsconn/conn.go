package wsconn

import (
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// Max time allowed to write a message to the client.
	writeWait = 5 * time.Second

	// Max time allowed to read the next pong message from client.
	pongWait = 30 * time.Second

	// Pings are sent to the client with this period.
	pingPeriod = 25 * time.Second
)

// Currently the connection is only used for sending text message i.e. (utf-8
// encoded messages)
type Conn struct {
	MessageType int // look at websocket.TextMessage / BinaryMessage

	*websocket.Conn
	reader     io.Reader // cached reader, which is reset whenever the messages is EOF
	writeWait  time.Duration
	pongWait   time.Duration
	pingPeriod time.Duration

	serverHeartbeat bool // if we sends 'h' character as heartbeat together with ping msg
}

func NewConnection(w *websocket.Conn) *Conn {
	c := &Conn{
		Conn:        w,
		pongWait:    pongWait,
		pingPeriod:  pingPeriod,
		writeWait:   writeWait,
		MessageType: websocket.TextMessage,
	}
	c.pingpong()
	return c
}

// NewConNewConnectionHeartbeat returns a connection with serverside heartbeat
// message support
//
// When serverHeartbeat is turned on, it sends 'h' frame as a standalone
// message together with the ping message
func NewConnectionHeartbeat(w *websocket.Conn, serverHeartbeat bool) *Conn {
	c := &Conn{
		Conn:            w,
		pongWait:        pongWait,
		pingPeriod:      pingPeriod,
		writeWait:       writeWait,
		serverHeartbeat: serverHeartbeat,
		MessageType:     websocket.TextMessage,
	}
	c.pingpong()
	return c
}

func (c *Conn) Write(b []byte) (int, error) {
	if c.writeWait != 0 {
		c.SetWriteDeadline(time.Now().Add(c.writeWait))
	}
	if err := c.WriteMessage(c.MessageType, b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *Conn) Read(b []byte) (int, error) {
	if c.reader == nil {
		t, r, err := c.NextReader()
		if err != nil {
			return 0, err
		}
		if t != c.MessageType {
			log.Printf("wsconn: ignored unknown message type. Got %d, wanted %d", t, c.MessageType)
			return 0, nil
		}
		c.reader = r
	}
	n, err := c.reader.Read(b)
	if err != nil {
		// When we have the frame EOF, treat it as no-op so we can receive the
		// next frame normally
		if err == io.EOF {
			c.reader = nil
			err = nil
		}
		return n, err
	}
	return n, nil
}

func (c *Conn) Close() error {
	c.SetWriteDeadline(time.Now().Add(c.writeWait))
	c.WriteMessage(websocket.CloseMessage, []byte{})
	return c.Conn.Close()
}

func (c *Conn) pingpong() {
	if c.pongWait != 0 {
		c.SetReadDeadline(time.Now().Add(c.pongWait))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(c.pongWait))
			return nil
		})
	}

	if c.pingPeriod != 0 {
		pingTicker := time.NewTicker(c.pingPeriod)
		go func() {
			defer pingTicker.Stop()
			for _ = range pingTicker.C {
				if c.writeWait != 0 {
					c.SetWriteDeadline(time.Now().Add(c.writeWait))
				}
				if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}

				// send serverHeartbeat if enabled
				if c.serverHeartbeat {
					c.Write([]byte("h"))
				}
			}
		}()
	}
}
