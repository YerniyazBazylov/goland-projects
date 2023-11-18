package data

import (
	"assignment3.ualikhan.net/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ClassicCars struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Year        int32     `json:"year,omitempty"`
	Cost        Cost      `json:"cost,omitempty"`
	Description string    `json:"description,omitempty"`
	Version     int32     `json:"version"`
}

func ValidateClassicCars(v *validator.Validator, classiccars *ClassicCars) {
	v.Check(classiccars.Name != "", "name", "must be provided")
	v.Check(len(classiccars.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(classiccars.Year != 0, "year", "must be provided")
	v.Check(classiccars.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(classiccars.Cost != 0, "cost", "must be provided")
	v.Check(classiccars.Cost > 0, "cost", "must be a positive integer")
}

type ClassicCarsModel struct {
	DB *sql.DB
}

func (m ClassicCarsModel) Insert(classiccars *ClassicCars) error {
	query := `
		INSERT INTO classic_cars (name, year, cost, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{classiccars.Name, classiccars.Year, classiccars.Cost, classiccars.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&classiccars.ID, &classiccars.CreatedAt, &classiccars.Version)
}

func (m ClassicCarsModel) Get(id int64) (*ClassicCars, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, year, cost, description, version
		FROM classic_cars
		WHERE id = $1`

	var classiccars ClassicCars

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&classiccars.ID,
		&classiccars.CreatedAt,
		&classiccars.Name,
		&classiccars.Year,
		&classiccars.Cost,
		&classiccars.Description,
		&classiccars.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &classiccars, nil
}

func (m ClassicCarsModel) Update(classiccars *ClassicCars) error {
	query := `
		UPDATE classic_cars
		SET name = $1, year = $2, cost = $3, description = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []interface{}{
		classiccars.Name,
		classiccars.Year,
		classiccars.Cost,
		classiccars.Description,
		classiccars.ID,
		classiccars.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&classiccars.Version)
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

func (m ClassicCarsModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM classic_cars
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

func (m ClassicCarsModel) GetAll(name string, filters Filters) ([]*ClassicCars, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT  count(*) OVER(), id, created_at, name, year, cost, description, version
		FROM classic_cars
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
	classiccars := []*ClassicCars{}

	for rows.Next() {
		var classiccar ClassicCars

		err := rows.Scan(
			&totalRecords,
			&classiccar.ID,
			&classiccar.CreatedAt,
			&classiccar.Name,
			&classiccar.Year,
			&classiccar.Cost,
			&classiccar.Description,
			&classiccar.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		classiccars = append(classiccars, &classiccar)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return classiccars, metadata, nil

}
