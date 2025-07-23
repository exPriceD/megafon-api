package repository

import "megafon-buisness-reports/internal/domain/entities"

func MapCityPhoneModelToEntity(m CityPhoneModel) entities.CityPhone {
	return entities.CityPhone{
		ID:             m.ID,
		City:           m.City,
		DiversionPhone: entities.Diversions(m.DiversionPhone),
	}
}
