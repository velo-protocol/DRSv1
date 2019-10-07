package main

import (
	"fmt"
	"gitlab.com/velo-labs/cen/libs/client2"
	vxdr "gitlab.com/velo-labs/cen/libs/xdr"
)

func main() {
	client := client2.DefaultTestNetClient
	client.VeloNodeURL = "localhost:8080"

	client.SourceAccount = "GBN4XRS4SGWOKAFOBTZ4JNZH7F3IJ3JUVDAFWD3ACES3KORW6QDAKGWA"
	client.SecretKey = "SAS2FS2L45PQY2Z4C4HAZWOP2GBTPN6Y6CL77Y6PCXTZMP5P2GIDWKH2"

	//client.SourceAccount = "GBN4XRS4SGWOKAFOBTZ4JNZH7F3IJ3JUVDAFWD3ACES3KORW6QDAKGWA" // for setupCredit
	//client.SecretKey = "SAS2FS2L45PQY2Z4C4HAZWOP2GBTPN6Y6CL77Y6PCXTZMP5P2GIDWKH2"     // for setupCredit

	client.Init()

	whiteListRequest := client2.WhitelistRequest{
		Address: "GBN4XRS4SGWOKAFOBTZ4JNZH7F3IJ3JUVDAFWD3ACES3KORW6QDAKGWA",
		Role:    vxdr.RoleTrustedPartner,
	}
	_, err := client.Whitelist(whiteListRequest)

	//setupCreditRequest := vclient.SetupCreditRequest{
	//	PeggedCurrency: "SGD",
	//	PeggedValue:    "1",
	//	AssetCode:      "vSGD",
	//}
	//_, err := client.SetupCredit(setupCreditRequest)

	//priceUpdateRequest := vclient.PriceUpdateRequest{
	//	Asset:                       "VELO",
	//	Currency:                    "THB",
	//	PriceInCurrencyPerAssetUnit: "1.1234567",
	//}
	//_, err := client.PriceUpdate(priceUpdateRequest)

	client.Close()

	fmt.Println(err)

}
