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

type User struct {
	Id           string `db:"id" json:"id"`
	Login        string `db:"login" json:"login"`
	PasswordHash string `db:"password_hash" json:"-"`
}
