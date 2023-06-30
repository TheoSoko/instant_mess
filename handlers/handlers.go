package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/TheoSoko/instant_mess/data"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var activeSockets = make(map[string]*websocket.Conn)
var sMutex sync.Mutex

type Payload struct {
	Message string `json:"message"`
}

/* (*-*) */

func Socketing(w http.ResponseWriter, r *http.Request) {
	var err error

	userId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("You need to provide and integer value ID as a query parameter"))
		return
	}

	// DISABLING AUTH FOR TESTING
	/*
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(401)
			w.Write([]byte("You need to provide a valid authentication token"))
			return
		}

		err = auth(token, strconv.Itoa(userId), w)
		if err != nil {
			// auth deals with http response messages
			w.WriteHeader(401)
			w.Write([]byte("You need to provide a valid authentication token"))
			return
		}
	*/

	socket, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer socket.Close()


	supplementaryID := uuid.New().String()
	connID := fmt.Sprint(userId) + "-" + supplementaryID
	sMutex.Lock()
	activeSockets[connID] = socket
	sMutex.Unlock()

	defer func() {
		sMutex.Lock()

		if as := activeSockets[connID]; as != nil {
			as.Close()
		}
		delete(activeSockets, connID)

		sMutex.Unlock()
	}()

	// For testing
	err = readFromSocket(socket)
	return
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

	// DISABLING AUTH FOR TESTING
		/*
			err = auth(token, id, w)
			if err != nil {
				// auth deals with http response
				return
			}
		*/

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
