package model

import (
	"database/sql"
	"errors"
)

var (
	ErrorRecordNotFound = errors.New("recored not found")
	ErrorEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
