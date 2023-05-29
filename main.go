package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}



func main(){

	messaging := func (w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r, nil)
		defer socket.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		
	}

	http.HandleFunc("/messaging", messaging)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}