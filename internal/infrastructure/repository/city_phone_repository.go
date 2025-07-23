package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"megafon-buisness-reports/internal/domain/entities"
	"megafon-buisness-reports/internal/domain/ports"
)

type CityPhoneRepository struct {
	db *sql.DB
}

func NewCityPhoneRepository(db *sql.DB) ports.CityPhoneRepository {
	return &CityPhoneRepository{db: db}
}

func (r *CityPhoneRepository) GetByCity(ctx context.Context, city string) (entities.CityPhone, error) {
	query := `SELECT city, ARRAY_AGG(diversion_phone ORDER BY diversion_phone) 
	FROM city_phones WHERE city = $1 GROUP BY city;
	`
	var m CityPhoneModel
	err := r.db.QueryRowContext(ctx, query, city).Scan(&m.City, &m.DiversionPhone)
	if errors.Is(err, sql.ErrNoRows) {
		return entities.CityPhone{}, fmt.Errorf("город %s не найден", city)
	}
	if err != nil {
		return entities.CityPhone{}, err
	}
	return MapCityPhoneModelToEntity(m), nil
}

func (r *CityPhoneRepository) GetCities(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT city FROM city_phones ORDER BY city;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var c string
		if err = rows.Scan(&c); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("справочник городов пуст")
	}
	return res, nil
}
