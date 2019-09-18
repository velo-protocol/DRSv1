package entities

import "time"

type Whitelist struct {
	ID   string
	StellarPublicAddress string
	RoleCode string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type WhitelistFilter struct {
	StellarPublicAddress *string
	RoleCode *string
}
