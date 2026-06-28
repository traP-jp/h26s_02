package repository

import "errors"

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrUniqueKeyDuplicated = errors.New("unique key duplicated")
	ErrViolatedForeignKey  = errors.New("foreign key violated")
)
