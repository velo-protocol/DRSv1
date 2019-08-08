package test_helpers

import "github.com/stellar/go/protocols/horizon"

func GetStellarAccount() horizon.Account {
	return horizon.Account{
		AccountID: "GD7R43KMK3AANO4TW722AKX6HZ7TKHKFZM5N4ASRUVU4FHB55V2JKOS2",
		Sequence: "396889337888770",
	}
}
