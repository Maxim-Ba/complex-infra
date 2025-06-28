package models

type UserCreate struct {
	Login    string
	Password string
}

type UserCreateRes struct {
	Login string
	Id    string
}
