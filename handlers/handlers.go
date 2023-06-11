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

func Socketing(w http.ResponseWriter, r *http.Request) {
	socket, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer socket.Close()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		socket.WriteMessage(5, []byte("You need to provide and integer value ID as a query parameter"))
		socket.Close()
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		socket.WriteMessage(5, []byte("You need to provide a valid authentication token"))
	}

	/* Auth here */
	req, _ := http.NewRequest("HEAD", "http://api.zemus.info:80/auth?id="+strconv.Itoa(id), nil)
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil  {
		log.Println("err at auth request :\n", err)
		socket.WriteMessage(5, []byte("An unknown error happened during authentication"))
		return
	}
	if res.StatusCode == 401 {
		socket.WriteMessage(5, []byte("The authentication failed"))
		return
	}
	if res.StatusCode != 204 {
		socket.WriteMessage(5, []byte("An unknown error happened during authentication"))
		return
	}
	/**/

	activeSockets[id] = socket
	defer delete(activeSockets, id)

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
	id, _ := mux.Vars(r)["id"]
	friendId, _ := mux.Vars(r)["friendId"]

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

	// ** Send the message to db
	err = data.PostMessage(intId, intFriendId, payload.Message)
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
		intId,
		" just successfully sent a message through a websocket : \n",
		payload.Message,
	)
	err = friendSocket.WriteMessage(1, []byte(message))
	if err != nil {
		// ** Send push notification if websocket fails.
		return
	}
}
