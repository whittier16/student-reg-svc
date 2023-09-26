package models

type User struct {
	email        string
	username     string
	passwordhash string
	fullname     string
	createDate   string
	role         int
}
