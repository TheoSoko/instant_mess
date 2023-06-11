package handlers

import (
	"encoding/json"
	"fmt"
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


func Socketing(w http.ResponseWriter, r *http.Request) {
	socket, err := wsUpgrader.Upgrade(w, r, nil)
	defer socket.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	id := r.URL.Query().Get("id")
	intId, err := strconv.Atoi(id)
	if id == "" || err != nil {
		socket.WriteMessage(5, []byte("You need to provide and integer value ID as a query parameter"))
		socket.Close()
		return
	}

	activeSockets[intId] = socket
	defer delete(activeSockets, intId)

	for {
		_, p, err := socket.ReadMessage()
		if err != nil {
			fmt.Println("erreur sur la lecture d'un message, ou fermeture du ws:", err)
			return
		}
		//fmt.Println("Message du client:", p)
		returnMessage := []byte(fmt.Sprint("We received your message ! It's : \"", string(p), "\""))
		socket.WriteMessage(1, returnMessage)
	}
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	friendId, ok2 := mux.Vars(r)["friendId"]
	if !ok || !ok2 {
		w.WriteHeader(400)
		w.Write([]byte("Please, provide the user id and the friend's id in the url"))
		return
	}

	intId, err := strconv.Atoi(id)
	intFriendId, err2 := strconv.Atoi(friendId)
	if err != nil || err2 != nil {
		w.WriteHeader(400)
		w.Write([]byte("Please, provide the user id and the friend's id as integer values"))
		return
	}

	_, err = data.GetUser(intFriendId)
	if err != nil {
		if err == fmt.Errorf("no_user") {
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

	// ** We would send the message to db here

	// ** Then, we write from the socket if friend is connected

	w.WriteHeader(204)

	friendSocket, ok := activeSockets[intFriendId]
	if !ok {
		// ** Send push notification if not connected
		return
	}

	message := fmt.Sprint("Hey, your friend with id ",
		intId,
		" just successfully sent a message through a websocket. The following : \n",
		payload.Message,
	)

	friendSocket.WriteMessage(1, []byte(message))


}
