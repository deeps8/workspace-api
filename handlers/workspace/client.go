package workspace

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"work-space-backend/database"
	"work-space-backend/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	// client id
	id string

	// hub that is connected to
	hub *Hub

	// websocket connection ref
	conn *websocket.Conn

	// channel for send the bytes back to UI
	send chan []byte
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024 * 1024 * 1024
)

/*
writePump pumps the messages from hub to websocket connection
goroutine is started for each connection
*/
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			{
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}

				w.Write(msg)

				n := len(c.send)
				for i := 0; i < n; i++ {
					w.Write([]byte{'\n'})
					w.Write(<-c.send)
				}

				if err := w.Close(); err != nil {
					return
				}
			}
		case <-ticker.C:
			{
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}

// readPump pumps data from websocket connection to hub
func (c *Client) readPump() {
	defer func() {
		// logMsg := &Message{ClientID: c.id, Text: fmt.Sprintf("User Left %s", c.id), MsgType: "info"}
		// c.hub.broadcast <- logMsg

		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, text, err := c.conn.ReadMessage()
		// log.Printf("(message type : %v)  value : %v", msgType, text)

		if err != nil {
			log.Printf("Readpump error : %v", err.Error())
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// 	log.Panicf("Connection closed with err : %v", err.Error())
			// }
			break
		}

		if string(text) == "con-closed" {
			c.hub.unregister <- c
			c.conn.Close()
			return
		}

		// msg := &Message{}
		// log.Printf("%v", text)
		// reader := bytes.NewReader(text)
		// decoder := json.NewDecoder(reader)
		// dErr := decoder.Decode(msg)
		// if dErr != nil {
		// 	log.Panicf("error while decoding msg : %v", dErr.Error())
		// }

		var brdData utils.RdbDataType = utils.RdbDataType{Data: string(text), Synced: false}
		brdJson, err := json.Marshal(brdData)
		if err != nil {
			log.Fatalf("Error marshaling struct to JSON: %v", err)
		}

		err = database.Rdb.Set(context.Background(), c.hub.roomId, string(brdJson), 0).Err()
		if err != nil {
			log.Fatalf("Error setting value in Redis: %v", err)
		}
		c.hub.broadcast <- &Message{Text: string(text), ClientID: c.id, MsgType: "msg"}
	}
}
