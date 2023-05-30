package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TheoSoko/instant_mess/data"
	"github.com/gorilla/mux"
)

func socketing(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	defer socket.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	id := r.URL.Query().Get("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		socket.WriteMessage(5, []byte("You need to provide and integer value ID as a query parameter"))
		socket.Close()
	}
	activeSockets[intID] = socket

	for {
		_, p, err := socket.ReadMessage()
		if err != nil {
			fmt.Println("erreur sur la lecture d'un messsage:", err)
			return
		}
		//fmt.Println("Message du client:", p)
		returnMessage := []byte(fmt.Sprint("We received your message ! It's : \"", string(p), "\""))
		socket.WriteMessage(1, returnMessage)
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	friendId, ok2 := mux.Vars(r)["friendId"]
	if !ok || !ok2 {
		w.WriteHeader(400)
		w.Write([]byte("Please, provide the user id and the friend's id in the url"))
		return
	}

	intID, err := strconv.Atoi(id)
	intFriendID, err := strconv.Atoi(friendId)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Please, provide the user id and the friend's id as integer values"))
		return
	}

	user := data.GetUser(intID)
	if !user.Friends[intFriendID] {
		w.WriteHeader(400)
		w.Write([]byte("Sorry, you can't send a message to a user you're not friends with"))
		return
	}

	var payload struct {
		Message string `json:"message"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	fmt.Println("payload from POST: ", payload.Message)

	// ** We would send the message to db here

	// ** Then, we write from the socket if friend is connected

	message := fmt.Sprint("Hey, your friend ",
		data.GetUser(intID).Firstname,
		" just successfully sent a message through a websocket. The following : \n",
		payload.Message,
	)
	if _, ok := activeSockets[intFriendID]; !ok {
		// ** Send push notification if not connected
		w.WriteHeader(204)
		return
	}

	activeSockets[intFriendID].WriteMessage(1, []byte(message))
}
