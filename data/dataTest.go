package data

/*
type User struct {
	Firstname string       `json:"firstname"`
	Lastname  string       `json:"lastname"`
	Friends   map[int]bool `json:"friends"`
}
*/

func GetUserTest(id int) User {

	users := map[int]User{
		1: {
			"John",
			"America",
			map[int]bool{2: true},
		},
		2: {
			"Rodrigo",
			"Perez",
			map[int]bool{1: true},
		},
	}

	return users[id]

}
