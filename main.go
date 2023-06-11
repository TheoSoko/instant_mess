package main

import (
	"fmt"
	"net/http"

	"github.com/TheoSoko/instant_mess/data"
	"github.com/TheoSoko/instant_mess/handlers"

	"github.com/gorilla/mux"
)

func main() {
	address := "127.0.0.1"
	port := "4000"
	r := mux.NewRouter()
	data.SqlConnect()

	http.Handle("/", r)
	r.HandleFunc("/ws", handlers.Socketing)
	r.HandleFunc("/users/{id}/friends/{friendId}/message", handlers.SendMessage).Methods("POST")

	err := http.ListenAndServe(address + ": " + port, r)
	if err != nil {
		panic(err)
	}


	fmt.Println("Listening on", address + ": " + port)
}
