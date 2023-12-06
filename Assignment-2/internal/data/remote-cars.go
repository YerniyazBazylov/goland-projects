package data

import (
	"assignment2.yerniyaz.net/internal/validator"
	"time"
)

type RemoteCars struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Year        int32     `json:"year,omitempty"`
	Cost        Cost      `json:"cost,omitempty"`
	Description string    `json:"description,omitempty"`
	Version     int32     `json:"version"`
}

func ValidateRemoteCars(v *validator.Validator, remotecars *RemoteCars) {
	v.Check(remotecars.Name != "", "name", "must be provided")
	v.Check(len(remotecars.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(remotecars.Year != 0, "year", "must be provided")
	v.Check(remotecars.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(remotecars.Cost != 0, "cost", "must be provided")
	v.Check(remotecars.Cost > 0, "cost", "must be a positive integer")
}
