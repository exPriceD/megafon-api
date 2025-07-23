package ports

import (
	"context"
	"megafon-buisness-reports/internal/domain/entities"
)

type CityPhoneRepository interface {
	GetByCity(ctx context.Context, city string) (entities.CityPhone, error)
	GetCities(ctx context.Context) ([]string, error)
}
