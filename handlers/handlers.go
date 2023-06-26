package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/TheoSoko/instant_mess/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var activeSockets = make(map[int]*websocket.Conn)
var sMutex sync.Mutex

type Payload struct {
	Message string `json:"message"`
}

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
		// authFromSocket deals with socket response messages
		socket.WriteMessage(8, []byte{0})
		socket.Close()
		return
	}

	sMutex.Lock()
	activeSockets[id] = socket
	sMutex.Unlock()

	defer func() {
		sMutex.Lock()

		if activeSockets[id] != nil {
			activeSockets[id].Close()
		}
		delete(activeSockets, id)

		sMutex.Unlock()
	}()

	// For testing
	readFromSocket(socket)
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

	var payload Payload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("A problem occured with the payload"))
	}

	// Send message to db
	err = data.PostMessage(id, friendId, payload.Message)
	if err != nil {
		w.WriteHeader(500)
	}

	// Message has been sent.
	w.WriteHeader(204)

	if err := writeToSocket(id, intFriendId, payload); err != nil {
		// (*-*) Send push notification if no active websocket for friend, or message failed.
	}

	return
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Ceci est le microservice de messagerie instantanée.\n" +
		"Allez à /ws pour ouvrir une connexion websocket, ou à : \n" +
		"\"POST /users/{id}/friends/{friendId}/message\" pour envoyer un message."))
	return
}
