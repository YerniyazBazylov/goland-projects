package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	RemoteCars  RemoteCarsModel
	Users       UserModel
	Permissions PermissionModel
	Tokens      TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		RemoteCars:  RemoteCarsModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
