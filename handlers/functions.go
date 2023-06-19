package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"log"

	"github.com/gorilla/websocket"
)

func authFromSocket(token string, id int, socket *websocket.Conn) error {
	req, _ := http.NewRequest("GET", "http://api.zemus.info/auth?id="+strconv.Itoa(id), nil)
	req.Header.Add("Authorization", token)
	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		socket.WriteMessage(1, []byte("An unknown error happened during authentication"))
		socket.WriteMessage(8, []byte{0})
		return err
	}
	if res.StatusCode == 401 {
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		socket.WriteMessage(1, []byte("The authentication failed, the response body : \n"+string(b)))
		socket.WriteMessage(8, []byte{0})
		return fmt.Errorf("401")
	}
	if res.StatusCode != 204 {
		socket.WriteMessage(1, []byte("An unknown error happened during authentication. Status from auth server :"+fmt.Sprint(res.StatusCode)))
		socket.WriteMessage(8, []byte{0})
		return fmt.Errorf("unknown")
	}

	return nil
}

func authFromMess(token string, id string, w http.ResponseWriter) error {
	req, _ := http.NewRequest("GET", "http://api.zemus.info/auth?id="+id, nil)
	req.Header.Add("Authorization", token)
	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("An unknown error happened during authentication"))
		return err
	}
	if res.StatusCode == 401 {
		w.WriteHeader(401)
		w.Write([]byte("The authentication failed"))
		log.Println("token: ", token)
		return fmt.Errorf("401")
	}
	if res.StatusCode != 204 {
		w.WriteHeader(401)
		w.Write([]byte("An unknown error happened during authentication. Status from auth server :" + fmt.Sprint(res.StatusCode)))
		return fmt.Errorf("unknown")
	}

	return nil
}
