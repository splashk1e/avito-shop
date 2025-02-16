package models

type User struct {
	Id        int
	Username  string
	PassHash  string
	Coins     int
	recieved  []Transaction
	sent      []Transaction
	Inventory []string
}
