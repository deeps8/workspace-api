package workspace

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"work-space-backend/database"
	"work-space-backend/utils"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var connUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		echo.New().Logger.Print(origin)
		// TODO: check for origins
		return true
	},
}

var hubList = NewHubList()

func ServeRoomWs(c echo.Context) error {
	roomid := c.Request().URL.Query().Get("roomid")
	userid := c.Request().URL.Query().Get("userid")
	if roomid == "" || userid == "" {
		log.Printf("Roomid: %v\nUserid: %v", roomid, userid)
		log.Fatal("Room or User ids are invalid")
		return c.JSON(http.StatusBadRequest, utils.Res{Message: "Room or User Id is not mentioned", Ok: false})
	}
	// upgrading http connection to websocket
	conn, wsErr := connUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if wsErr != nil {
		log.Printf("Error while upgrading http connection to websocket")
	}

	/*
		1. if roomid exists or not
		2. if not create one and register the client
		3. if yes register the client in that room
		4. get the "data" from redis if not then fetch from db
	*/
	hub, ok := hubList.Hubs[roomid]
	log.Printf("Roomid : %v", hub)
	if !ok {
		log.Printf("Room Does not exists")
		// room does not exists, create one
		hub = newHub(roomid)
		hubList.Hubs[roomid] = hub
		go hub.run()
	}

	var brdData utils.RdbDataType
	rdata, rerr := database.Rdb.Get(context.Background(), roomid).Result()

	if rerr != nil || rdata == "" {
		brd, berr := database.SelectBordData(roomid)
		if berr != nil {
			return c.JSON(http.StatusNotFound, utils.Res{Message: berr.Error(), Ok: false})
		}
		brdData.Data = brd.Data
		brdData.Synced = true
		brdJson, err := json.Marshal(brdData)
		if err != nil {
			log.Fatalf("Error marshaling struct to JSON: %v", err)
		}

		err = database.Rdb.Set(context.Background(), roomid, string(brdJson), 0).Err()
		if err != nil {
			log.Fatalf("Error setting value in Redis: %v", err)
		}
	} else {
		err := json.Unmarshal([]byte(rdata), &brdData)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON from Redis: %v", err)
		}
	}

	hub.messages = &Message{ClientID: userid, Text: brdData.Data, MsgType: "msg"}
	// creating a new client for each time use joins
	newClient := &Client{id: userid, hub: hub, conn: conn, send: make(chan []byte, 512)}

	// registering the newly created client to hub
	newClient.hub.register <- newClient

	// now start the goroutines for read and write the data
	go newClient.writePump()
	go newClient.readPump()

	// logMsg := &Message{ClientID: newClient.id, Text: fmt.Sprintf("User Joined %s", newClient.id), MsgType: "info"}
	// newClient.hub.broadcast <- logMsg
	return nil
}
