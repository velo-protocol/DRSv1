package entities

type WhiteList struct {
	ID               string
	StellarPublicKey string
	RoleCode         string
}

type WhiteListFilter struct {
	StellarPublicKey *string
	RoleCode         *string
}
