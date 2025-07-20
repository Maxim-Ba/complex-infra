package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"go-auth/internal/app"
	"go-auth/internal/models"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage( db app.DB) *UserStorage {
	fmt.Println(db)
	return &UserStorage{db: db.GetConnection()}
}

func (s *UserStorage) Save(user models.UserCreateDto) (models.UserCreateRes, error) {

var res models.UserCreateRes
	


	query := `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
		RETURNING id, login
	`

	err := s.db.QueryRow(query, user.Login,  user.PasswordHash).Scan(&res.Id, &res.Login)
	if err != nil {
		return res, fmt.Errorf("failed to save user: %w", err)
	}

	return res, nil
}

func (s *UserStorage) Get(user models.UserCreateDto) (*models.UserCreateRes, error) {
	var res models.UserCreateRes

	query := `
		SELECT id, login
		FROM users
		WHERE login = $1 AND password_hash = $2
	`

	err := s.db.QueryRow(query, user.Login, user.PasswordHash).Scan(&res.Id, &res.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &res, nil
}

func (s *UserStorage) Update(user models.UserCreateDto) error {
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE login = $2
	`

	result, err := s.db.Exec(query, user.PasswordHash, user.Login)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
func (s *UserStorage) GetById(userId string) (*models.UserCreateRes, error) {
var res models.UserCreateRes

	query := `
		SELECT id
		FROM users
		WHERE password_hash = $1
	`

	err := s.db.QueryRow(query, userId).Scan(&res.Id, &res.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &res, nil
}
