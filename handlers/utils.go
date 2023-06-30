package handlers

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

func writeToSocket(id string, friendId int, payload Payload) error {
	var err error
	friendSockets := findSockets(friendId)

	if len(friendSockets) > 0 {
		// We write to the socket if friend is connected.
		message := fmt.Sprint("Your friend ", id, " wrote : \n", payload.Message)
		for _, ws := range friendSockets {
			sMutex.Lock()
			err = ws.WriteMessage(1, []byte(message))
			sMutex.Unlock()
		}
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("no_conn")
}

func readFromSocket(socket *websocket.Conn) {
	for {
		// THIS IS FOR TESTING, the socket won't be accessed from here, so no conflict with goroutines.
		_, p, err := socket.ReadMessage()
		if err != nil {
			fmt.Println("erreur sur la lecture d'un message, ou fermeture du ws:", err)
			return
		}
		returnMessage := []byte(fmt.Sprint("We received your message ! It's : \"", string(p), "\""))
		socket.WriteMessage(1, returnMessage)
	}

}

func findSockets(friendId int) []*websocket.Conn {
	wsConnBuff := []*websocket.Conn{}
	for id, ws := range activeSockets {
		if strings.Split(id, "-")[0] == fmt.Sprint(friendId) {
			wsConnBuff = append(wsConnBuff, ws)
		}
	}

	return wsConnBuff
}
