package models

type UserCreateReq struct {
	Login    string
	Password string
}

type UserCreateRes struct {
	Login string
	Id    string
}

type UserCreateDto struct {
	Login        string
	PasswordHash string
}
