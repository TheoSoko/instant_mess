package data

import (
	"database/sql"
	"fmt"
)

type User struct {
	Firstname string       `json:"firstname"`
	Lastname  string       `json:"lastname"`
	Friends   map[int]bool `json:"friends"`
}

func GetUser(id int) (User, error) {
	var user User

	row := db.QueryRow("SELECT `firstname`, `lastname` FROM users WHERE id = ?", id)
	if err := row.Scan(&user.Firstname, &user.Lastname); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("no_user")
		}
		fmt.Println("error at GetUser", id, err)
		return user, err
	}

	return user, nil
}

func PostMessage(senderId int, receiverId int, message string) error {
	_, err := db.Exec(
		"INSERT INTO `messages` (user_sender_id, user_receiver_id, content) VALUES (?, ?, ?)",
		senderId, receiverId, message)
	if err != nil {
		return err
	}
	return nil
}
