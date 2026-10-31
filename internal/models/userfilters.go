package models

type UsersFiltres struct {
	Active  bool
	Limit   int
	Offset  int
	Role    string
	Country string
	OrderBy string
}
