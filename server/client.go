package server

import (
	// Utils
	"fmt"

	// Gorilla websockets
	"github.com/gorilla/websocket"
	"github.com/hexagon/gotris/util"

	// Redis
	"gopkg.in/mgo.v2"
)

func Client(c *websocket.Conn, r *mgo.Session) {

	// Set connection options
	wsOutChannel := make(chan string, 3)
	wsInChannel := make(chan map[string]interface{}, 3)

	defer func() {
		close(wsInChannel)
	}()

	go Player(wsInChannel, wsOutChannel, r)
	go wsWriter(c, wsOutChannel)

	// Incoming websocket messages is received and parsed using the wsReader goroutine
	// which pass the unmarshalled data to wsInChannel

	wsReader(c, wsInChannel)

}

func wsReader(conn *websocket.Conn, wsInChannel chan map[string]interface{}) {

	for {

		// Read the actual message
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Unmarshal JSON into a map
		unmarshalErr, packet := util.Unmarshal(message)
		if unmarshalErr != nil {
			fmt.Println("Error unmarshaling incoming JSON: %v", unmarshalErr)
		}

		// Send received object to player
		wsInChannel <- packet

	}
}

func wsWriter(conn *websocket.Conn, wsOutChannel chan string) {

	// Write pump
	for {
		select {
		case data, ok := <-wsOutChannel:
			if ok {
				err := conn.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					return
				}
			} else {
				return
			}
		}
	}

}
