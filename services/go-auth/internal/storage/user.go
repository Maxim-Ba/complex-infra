package storage

import "go-auth/internal/models"

type UserStorage struct {
}

func NewUserStorage() *UserStorage {
	return &UserStorage{}
}

func (s *UserStorage) Save(user models.UserCreateReq) models.UserCreateRes {
var u models.UserCreateRes
	return u
}

func (s *UserStorage) Get(user models.UserCreateDto) (*models.UserCreateRes, error) {
var u models.UserCreateRes
	return &u, nil
}
