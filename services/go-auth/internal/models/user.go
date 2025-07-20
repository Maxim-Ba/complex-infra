package models

// UserCreateReq represents user registration/login request
// @Description Данные для регистрации или входа пользователя
type UserCreateReq struct {
	Login    string `json:"login" example:"user123" minLength:"3" maxLength:"20"`
	Password string `json:"password" example:"strongPassword123" minLength:"6" maxLength:"32"`
}

// UserCreateRes represents successful registration response
// @Description Ответ при успешной регистрации пользователя
type UserCreateRes struct {
	Login string `json:"login" example:"user123"`
	Id    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
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
