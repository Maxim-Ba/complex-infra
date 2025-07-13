package services

import "errors"

var ErrUserExists = errors.New("user allready exists")
var ErrWrongLoginOrPassword = errors.New("wrong password or login")
var ErrLoginAndPasswordAreRequired = errors.New("login and password are required")
