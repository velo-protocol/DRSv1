package vxdr

// Role is a constant which defined supported role.
type Role string

const (
	RoleTrustedPartner Role = "TRUSTED_PARTNER"
	RolePriceFeeder    Role = "PRICE_FEEDER"
	RoleRegulator      Role = "REGULATOR"
)

// IsValid checks if the given role is supported/valid or not.
func (role Role) IsValid() bool {
	return role == RolePriceFeeder ||
		role == RoleTrustedPartner ||
		role == RoleRegulator
}

// RoleMap defines a pretty string for each Role constant.
var RoleMap = map[Role]string{
	RoleTrustedPartner: "Trusted Partner",
	RolePriceFeeder:    "Price Feeder",
	RoleRegulator:      "Regulator",
}
