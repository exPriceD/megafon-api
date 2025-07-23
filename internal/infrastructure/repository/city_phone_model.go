package repository

import "github.com/lib/pq"

type CityPhoneModel struct {
	ID             int64          `db:"id"`
	City           string         `db:"city"`
	DiversionPhone pq.StringArray `db:"diversion_phone"`
}
