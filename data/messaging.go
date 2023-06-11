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
			fmt.Println("GetUser", id, ": pas d'user avec cet id")
			return user, fmt.Errorf("no_user")
		}
		fmt.Println("GetUser", id, err)
		return user, fmt.Errorf("unknown")
	}

	return user, nil
}
