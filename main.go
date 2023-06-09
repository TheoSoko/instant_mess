package main

import (
	"net/http"

	"github.com/TheoSoko/instant_mess/handlers"
	"github.com/gorilla/mux"
)


func main() {

	r := mux.NewRouter()

	http.Handle("/", r)
	r.HandleFunc("/ws", handlers.Socketing)
	r.HandleFunc("/users/{id}/friends/{friendId}/message", handlers.SendMessage).Methods("POST")

	err := http.ListenAndServe("127.0.0.1:4000", r)
	if err != nil {
		panic(err)
	}
}
