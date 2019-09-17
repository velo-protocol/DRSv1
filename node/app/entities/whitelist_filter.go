package entities

// WhitelistFilter for filtering
type WhitelistFilter struct {
	StellarAddress *string `form:"stellar_address,omitempty"`
	Role *string `form:"role,omitempty"`
}
