package vxdr

type Role string

const (
	RoleTrustedPartner Role = "TRUSTED_PARTNER"
	RolePriceFeeder    Role = "PRICE_FEEDER"
	RoleRegulator      Role = "REGULATOR"
)

func (role Role) IsValid() bool {
	return role == RolePriceFeeder ||
		role == RoleTrustedPartner ||
		role == RoleRegulator
}

var RoleMap = map[Role]string{
	RoleTrustedPartner: "Trusted Partner",
	RolePriceFeeder:    "Price Feeder",
	RoleRegulator:      "Regulator",
}
