package vclient

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"log"
)

func Example() {
	client, err := NewDefaultTestNetClient("localhost:8080", clientSecretKey)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	txResult, err := client.WhiteList(context.Background(), vtxnbuild.WhiteList{
		Address:  "GC5F4E7IKMDFNOL7Z5WDHC42LBLLQL2UFY6KQALO2RRHC5EMJJRECPI3",
		Role:     "PRICE_FEEDER",
		Currency: "USD",
	})
	if err != nil {
		panic(err)
	}
	log.Println(txResult.TransactionSuccessToString())

	txResult, err = client.SetupCredit(context.Background(), vtxnbuild.SetupCredit{
		PeggedValue:    "1.0",
		PeggedCurrency: "USD",
		AssetCode:      "vUSD",
	})
	if err != nil {
		panic(err)
	}
	log.Println(txResult.TransactionSuccessToString())

	txResult, err = client.PriceUpdate(context.Background(), vtxnbuild.PriceUpdate{
		Asset:                       "VELO",
		Currency:                    "USD",
		PriceInCurrencyPerAssetUnit: "0.5",
	})
	if err != nil {
		panic(err)
	}
	log.Println(txResult.TransactionSuccessToString())
}
