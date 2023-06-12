package main

import (
	"log"
	"net/http"

	"github.com/TheoSoko/instant_mess/data"
	"github.com/TheoSoko/instant_mess/handlers"

	"github.com/gorilla/mux"
)

func main() {
	address := "0.0.0.0"
	port := "4000"
	data.SqlConn()

	r := mux.NewRouter()
	r.Handle("/", r)
	r.HandleFunc("/ws", handlers.Socketing)
	r.HandleFunc("/users/{id}/friends/{friendId}/message", handlers.SendMessage).Methods("POST")

	log.Println("Listening on", address+": "+port)
	err := http.ListenAndServe(address+": "+port, r)
	if err != nil {
		panic(err)
	}
}
