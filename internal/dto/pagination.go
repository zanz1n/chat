package dto

import "github.com/google/uuid"

type Pagination struct {
	LastSeen uuid.UUID `json:"last_seen"`
	//govalid:gt=0
	//govalid:lte=100
	Limit int `json:"limit"`
}
