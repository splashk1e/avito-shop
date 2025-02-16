package models

type Transaction struct {
	Id       int
	Reciever string
	Sender   string
	Amount   int
	Item     string
}
