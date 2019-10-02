package entities

type Whitelist struct {
	ID               string
	StellarPublicKey string
	RoleCode         string
}

type WhitelistFilter struct {
	StellarPublicKey *string
	RoleCode         *string
}
