package main

import (
	"log"
	"net/http"

	"github.com/TheoSoko/instant_mess/data"
	"github.com/TheoSoko/instant_mess/handlers"

	"os"

	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)

	log.Println("file: ", b)
	log.Println("Directory: ", basepath)

	err = godotenv.Load(basepath + "/env/.env")
	if err != nil {
		log.Fatal(err)
	}

	err = data.SqlConn()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Hello)
	r.HandleFunc("/ws", handlers.Socketing)
	r.HandleFunc("/users/{id}/friends/{friendId}/message", handlers.SendMessage).Methods("POST")

	log.Println("Listening on", os.Getenv("IP")+": "+os.Getenv("PORT"))
	err = http.ListenAndServe(os.Getenv("IP")+":"+os.Getenv("PORT"), r)
	if err != nil {
		panic(err)
	}
}
