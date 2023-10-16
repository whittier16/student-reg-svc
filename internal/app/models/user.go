package models

type User struct {
	Email        string
	Username     string
	Passwordhash string
	Fullname     string
	CreateDate   string
	Role         int
}
