package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var activeSockets = make(map[int]*websocket.Conn)

func main() {

	r := mux.NewRouter()

	http.Handle("/", r)
	r.HandleFunc("/ws", socketing)
	r.HandleFunc("/users/{id}/friends/{friendId}/message", sendMessage).Methods("POST")

	err := http.ListenAndServe("127.0.0.1:6969", r)
	if err != nil {
		panic(err)
	}
}
