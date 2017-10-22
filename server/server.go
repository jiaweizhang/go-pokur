package server

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
)

var tm map[string](map[int64]*websocket.Conn)

type message struct {
	Player string `json:"player"`
	MType  string `json:"mType"`
	Msg    string `json:"msg"`
}

// initializes server
func Init() {
	tm = make(map[string](map[int64]*websocket.Conn))

	// setup /table/:tableName
	setupTable("a")

}

func tableHandler(tableName string) websocket.Handler {

	tf := func(ws *websocket.Conn) {

		var err error

		for {
			var received message
			if err = websocket.JSON.Receive(ws, &received); err != nil {
				log.Printf("Can't receive: %v", err)
				break
			}

			// attempt to parse player
			playerID, err := strconv.ParseInt(received.Player, 10, 64)
			if err != nil {
				log.Printf("failed to parse player id: %v", err)
				if err = websocket.Message.Send(ws, "Cannot parse player ID"); err != nil {
					log.Println("Can't send parse player ID")
					break
				}
			}

			// TODO(jiaweizhang) check player membership status

			// if mtype == new, check to make sure player is not already connected
			if received.MType == "new" {
				if _, ok := tm[tableName][playerID]; ok {
					// already exists, so deny request
					if err = websocket.Message.Send(ws, "Player is already connected"); err != nil {
						log.Println("Cannot send player already connected")
						break
					}
				} else {
					// add player's ws to map
					tm[tableName][playerID] = ws
					if err = websocket.Message.Send(ws, "Player successfully connected"); err != nil {
						log.Println("Cannot send player successfully connected")
						break
					}
				}
			} else if received.MType == "msg" {
				// receiving message
				log.Printf("Received msg from player %s: %s", received.Player, received.Msg)
				// TODO(jiaweizhang) channel integration
				// send message to everyone
				for _, wsws := range tm[tableName] {
					if err = websocket.Message.Send(wsws, received.Player+": "+received.Msg); err != nil {
						log.Printf("Cannot send message to a player: %v", err)
						break
					}
				}
			}
		}
	}

	return websocket.Handler(tf)
}

func setupTable(tableName string) error {
	if _, ok := tm[tableName]; !ok {
		// create map from player id to their websocket
		tm[tableName] = make(map[int64]*websocket.Conn)
	}
	http.Handle("/table/"+tableName, tableHandler(tableName))
	return nil
}
