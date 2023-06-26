package handlers

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func writeToSocket(id string, friendId int, payload Payload) error {
	friendSocket, exists := activeSockets[friendId]
	
	if exists {
		// We write to the socket if friend is connected.
		message := fmt.Sprint("Hey, your friend with id ", id,
			" just successfully sent a message through a websocket : \n",
			payload.Message,
		)
		sMutex.Lock()
		err := friendSocket.WriteMessage(1, []byte(message))
		sMutex.Unlock()
		if err == nil {
			return err
		}
	}

	return fmt.Errorf("no_conn")
}

func readFromSocket(socket *websocket.Conn){
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
