package data

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	Firstname string       `json:"firstname"`
	Lastname  string       `json:"lastname"`
	Friends   map[int]bool `json:"friends"`
}

func GetUser(id string) (User, error) {
	var user User

	if db == nil {
		fmt.Println("db == nil")
		return user, nil
	}

	row := db.QueryRow("SELECT `firstname`, `lastname` FROM users WHERE id = ?", id)
	if err := row.Scan(&user.Firstname, &user.Lastname); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("no_user")
		}
		log.Println("error at GetUser", id, err)
		return user, err
	}

	return user, nil
}

func PostMessage(senderId string, receiverId string, message string) error {
	_, err := db.Exec(
		"INSERT INTO `messages` (user_sender_id, user_receiver_id, content) VALUES (?, ?, ?)",
		senderId, receiverId, message)
	if err != nil {
		return err
	}
	return nil
}
