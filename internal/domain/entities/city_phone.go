package entities

type Diversions []string

type CityPhone struct {
	ID             int64
	City           string
	DiversionPhone Diversions
}
