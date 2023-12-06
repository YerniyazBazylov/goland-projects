package data

import (
	"assignment3.yerniaz.net/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
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

type RemoteCarsModel struct {
	DB *sql.DB
}

func (m RemoteCarsModel) Insert(remotecars *RemoteCars) error {
	query := `
		INSERT INTO remote_cars (name, year, cost, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{remotecars.Name, remotecars.Year, remotecars.Cost, remotecars.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&remotecars.ID, &remotecars.CreatedAt, &remotecars.Version)
}

func (m RemoteCarsModel) Get(id int64) (*RemoteCars, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, year, cost, description, version
		FROM remote_cars
		WHERE id = $1`

	var remotecars RemoteCars

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&remotecars.ID,
		&remotecars.CreatedAt,
		&remotecars.Name,
		&remotecars.Year,
		&remotecars.Cost,
		&remotecars.Description,
		&remotecars.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &remotecars, nil
}

func (m RemoteCarsModel) Update(remotecars *RemoteCars) error {
	query := `
		UPDATE remote_cars
		SET name = $1, year = $2, cost = $3, description = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []interface{}{
		remotecars.Name,
		remotecars.Year,
		remotecars.Cost,
		remotecars.Description,
		remotecars.ID,
		remotecars.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&remotecars.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m RemoteCarsModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM remote_cars
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m RemoteCarsModel) GetAll(name string, filters Filters) ([]*RemoteCars, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT  count(*) OVER(), id, created_at, name, year, cost, description, version
		FROM remote_cars
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	remotecars := []*RemoteCars{}

	for rows.Next() {
		var remotecar RemoteCars

		err := rows.Scan(
			&totalRecords,
			&remotecar.ID,
			&remotecar.CreatedAt,
			&remotecar.Name,
			&remotecar.Year,
			&remotecar.Cost,
			&remotecar.Description,
			&remotecar.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		remotecars = append(remotecars, &remotecar)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return remotecars, metadata, nil

}
