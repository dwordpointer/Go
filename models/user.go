package models

type User struct {
	Id        uint
	Firstname string
	LastName  string
	Email     string `gorm:unique`
	Password  []byte
}
