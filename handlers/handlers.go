package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TheoSoko/instant_mess/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var activeSockets = make(map[int]*websocket.Conn)

/* (*-*) */

func Socketing(w http.ResponseWriter, r *http.Request) {
	socket, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		socket.WriteMessage(1, []byte("You need to provide and integer value ID as a query parameter"))
		socket.WriteMessage(8, []byte{0})
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		socket.WriteMessage(1, []byte("You need to provide a valid authentication token"))
		socket.WriteMessage(8, []byte{0})
		return
	}

	err = authFromSocket(token, id, socket)
	if err != nil {
		// authFromSocket deals with socket response
		return
	}

	activeSockets[id] = socket
	defer delete(activeSockets, id)

	for {
		_, p, err := socket.ReadMessage()
		if err != nil {
			fmt.Println("erreur sur la lecture d'un message, ou fermeture du ws:", err)
		}
		//fmt.Println("Message du client:", p)
		returnMessage := []byte(fmt.Sprint("We received your message ! It's : \"", string(p), "\""))
		socket.WriteMessage(1, returnMessage)
	}
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	friendId, _ := mux.Vars(r)["friendId"]

	_, err := strconv.Atoi(id)
	intFriendId, err2 := strconv.Atoi(friendId)
	if err != nil || err2 != nil {
		w.WriteHeader(400)
		w.Write([]byte("Please, provide the user id and the friend's id as integer values"))
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(401)
		w.Write([]byte("You need to provide a bearer authorization token"))
		return
	}

	err = authFromMess(token, id, w)
	if err != nil {
		// authFromMess deals with http response
		return
	}

	_, err = data.GetUser(friendId)
	if err != nil {
		if err.Error() == "no_user" {
			w.WriteHeader(404)
			w.Write([]byte("No user is associated with that id"))
			return
		}
		w.WriteHeader(500)
		return
	}

	var payload struct {
		Message string `json:"message"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("A problem occured with the payload"))
	}

	// ** Send the message to db
	err = data.PostMessage(id, friendId, payload.Message)
	if err != nil {
		w.WriteHeader(500)
	}

	// Message has been sent.
	w.WriteHeader(204)

	friendSocket, ok := activeSockets[intFriendId]

	if !ok {
		// ** Send push notification if not connected.
		return
	}

	// ** We write to the socket if friend is connected.
	message := fmt.Sprint("Hey, your friend with id ",
		id,
		" just successfully sent a message through a websocket : \n",
		payload.Message,
	)
	err = friendSocket.WriteMessage(1, []byte(message))
	if err != nil {
		// ** Send push notification if websocket fails.
		return
	}
}
