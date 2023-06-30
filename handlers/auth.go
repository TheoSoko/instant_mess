package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

func authFromSocket(token string, id int, socket *websocket.Conn) error {
	protocol := os.Getenv("AUTH_PROTOCOL")
	ip := os.Getenv("AUTH_IPV6")
	port := os.Getenv("AUTH_PORT")
	strId := strconv.Itoa(id)

	req, err := http.NewRequest("GET", protocol+"://"+ip+":"+port+"/auth?id="+strId, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", token)
	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		socket.WriteMessage(1, []byte("An unknown error happened during authentication"))
		return err
	}
	if res.StatusCode == 401 {
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		socket.WriteMessage(1, []byte("The authentication failed, the response body : \n"+string(b)))
		return fmt.Errorf("401")
	}
	if res.StatusCode != 204 {
		socket.WriteMessage(1, []byte("An unknown error happened during authentication. Status from auth server :"+fmt.Sprint(res.StatusCode)))
		return fmt.Errorf("unknown")
	}

	return nil
}

func auth(token string, id string, w http.ResponseWriter) error {
	protocol := os.Getenv("AUTH_PROTOCOL")
	ip := os.Getenv("AUTH_IPV6")
	port := os.Getenv("AUTH_PORT")

	req, _ := http.NewRequest("GET", protocol+"://"+ip+":"+port+"/auth?id="+id, nil)

	req.Header.Add("Authorization", token)
	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("An unknown error happened during authentication"))
		fmt.Println("auth erreur : ",err)
		return err
	}
	if res.StatusCode == 401 {
		//defer res.Body.Close()
		//b, _ := io.ReadAll(res.Body)
		w.WriteHeader(401)
		w.Write([]byte("The authentication failed"))
		return fmt.Errorf("401")
	}
	if res.StatusCode != 204 {
		w.WriteHeader(401)
		w.Write([]byte("An unknown error happened during authentication. Status from auth server :" + fmt.Sprint(res.StatusCode)))
		return fmt.Errorf("unknown")
	}

	return nil
}
