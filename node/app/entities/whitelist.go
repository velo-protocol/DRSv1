package entities

type Whitelist struct {
	ID   string
	StellarAddress string
	Role string
}

type WhitelistFilter struct {
	StellarAddress *string
	Role *string
}
