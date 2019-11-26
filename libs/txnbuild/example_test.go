package vtxnbuild

import (
	"github.com/stellar/go/txnbuild"
	"log"
)

func ExampleMintCredit() {
	veloTxB64, err := (&VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &MintCredit{
			AssetCodeToBeIssued: "<ASSET CODE FOR ISSUED>",        // Ex: vTHB
			CollateralAssetCode: "<COLLATERAL ASSET CODE>",        // Ex: VELO (Now Supported Only VELO Token)
			CollateralAmount:    "<COLLATERAL AMOUNT FOR ISSUED>", // Ex: 10
		},
	}).BuildSignEncode(kp1, kp2)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(veloTxB64)
	// Output: AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAwAAAAAAAAAAAAAAAAAAAAEAAAAEdlRIQgAAAARWRUxPAAAAAAX14QAAAAAAAAAAAAAAAALW0+y7AAAAQFaolJTJnRVzsfuFp60qMqdfU4qX6rIvbuWfkPKYmM339ELqOydwraLaMylPG+wTFjTi/9YkVhdzyDKNK5KsuAPkR85aAAAAQF3h6oYQ8PuR5RxY34xLPsWWbm3T3VOUQNxbgXg0r3jaJCuVnoEKE4VaO5JkntHEXTyf+ookg/2kAg2p5XNj5QI=
}

func ExamplePriceUpdate() {
	veloTxB64, err := (&VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &PriceUpdate{
			Asset:                       "<COLLATERAL ASSET CODE>",            // Ex: VELO (Now Supported Only VELO Token)
			Currency:                    "<FIAT CURRENCY>",                    // Ex: THB (Now supported Only THB, USD and SGD)
			PriceInCurrencyPerAssetUnit: "<PRICE IN CURRENCY PER ASSET UNIT>", // Ex: 20
		},
	}).BuildSignEncode(kp1, kp2)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(veloTxB64)
	// Output: AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAgAAAAAAAAAAAAAAAQAAAARWRUxPAAAAA1RIQgAAAAAAC+vCAAAAAAAAAAAAAAAAAAAAAALW0+y7AAAAQIwiy2a9e3zHFnyYaetdYEF0uRj2xYkqtc2tBVSgzJRo2jA7k9eJ7O61AQo8HCMt2m9Y2gLaVJQuplws1VOOTA7kR85aAAAAQEV5p/hzjgXv4bVK38YDgliRCeqny0wQc0OOXhmosdZXjQ6KGUBmrQJOEC2stObkFclVTr0/OxcmFQumCCY3AwM=
}

func ExampleRebalanceReserve() {
	veloTxB64, err := (&VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &RebalanceReserve{},
	}).BuildSignEncode(kp1, kp2)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(veloTxB64)
	// Output: AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAALW0+y7AAAAQF02xQ66AyeS6jes18Pvuz3ZADLch3Li60sQpVxAAc3IhNEjveBs7K/U0w7MDRw054lXnIPFROuzXovbz3+C9wPkR85aAAAAQOHQLGY/p0BpxkCppGwvmeSXY0eQsjsRdBO1fGiV1hHcLzT/jexCb3SIIEm7Eo1ragq6GTNrY14pb4+9XZj03wA=
}

func ExampleRedeemCredit() {
	veloTxB64, err := (&VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &RedeemCredit{
			AssetCode: "<STABLE CREDIT ASSET CODE>",      // Ex: vTHB
			Issuer:    "<STABLE CREDIT ISSUER ADDRESS>",  // Ex: GAXKPU22AE22NO7FXSW7GTNJJ6FGN5NQLXWTJGNBF4VOKLXVJ3RROXTI
			Amount:    "<STABLE CREDIT REDEEMED AMOUNT>", // Ex: 100 (Amount of vTHB)
		},
	}).BuildSignEncode(kp1, kp2)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(veloTxB64)
	// Output: AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAABAAAAAAAAAAAAAAAAAAAAAAAAAABAAAABHZUSEIAAAAALqfTWgE1prvlvK3zTalPimb1sF3tNJmhLyrlLvVO4xcAAAAAO5rKAAAAAAAAAAAC1tPsuwAAAECIrN/vhzms66OV4EFq2NcY9JWc5eNq/vNuiC5VV5+db9zBnFH+4uv0FLKoMUuVH+xdRl6jNSdQH3XTKAn2BdkO5EfOWgAAAECMgULpNumjGDLgevboF7206fm8xa+dUe5JVlv5Bbzq1g8uWDvbJdUaF0OCqR9qbDzpqRErMa9MK+g+R1QhgIgB
}

func ExampleSetupCredit() {
	veloTxB64, err := (&VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &SetupCredit{
			PeggedValue:    "<PEGGED VALUE OF STABLE CREDIT>",    // Ex: 1
			PeggedCurrency: "<PEGGED CURRENCY OF STABLE CREDIT>", // Ex: THB
			AssetCode:      "<ASSET CODE OF STABLE CREDIT>",      // Ex: vTHB
		},
	}).BuildSignEncode(kp1, kp2)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(veloTxB64)
	// Output: AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAQAAAAAAAAABAAAAAACYloAAAAADVEhCAAAAAAR2VEhCAAAAAAAAAAAAAAAAAAAAAAAAAALW0+y7AAAAQFZqHpPLIbA5wSwEybqbtgApla5PKJRI/q90DoWO/AEbyF2QNRfUnuYcNfxgK8PKvbCNycUedtWRA0MwSGY9BwLkR85aAAAAQBQ2OOtzW2XHzNyK96FKLm2Q3ri6yqm+qW3vzvqcH/2GEbmBTDh+divAUjVywMgxAC+OhoucicyfqbGaRf6GNAw=
}

func ExampleWhitelist() {
	veloTxB64, err := (&VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &Whitelist{
			Address: "<STELLAR ACCOUNT ADDRESS>", // Ex: GC2ROYZQH5FTVEPQZF7CAB32SCJC7DWVKILDUAT5BCU5O7HEI7HFUB25
			Role:    "<VELO ROLE>",               // EX: TRUSTED_PARTNER  (Now supported Only TRUSTED_PARTNER, PRICE_FEEDER and REGULATOR)
		},
	}).BuildSignEncode(kp1, kp2)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(veloTxB64)
	// Output: AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAAAAAAEAAAAAtRdjMD9LOpHwyX4gB3qQki+O1VIWOgJ9CKnXfORHzloAAAAPVFJVU1RFRF9QQVJUTkVSAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALW0+y7AAAAQEQI7/oYaybfF0m5jQDO72fFpzx1YyyiJn7keTGp34B6hpgbCobN8M41WERPEP+Z+kRSxDkzJe49GT7jicQhTwLkR85aAAAAQFKPYrtXStZiH4mDSIcmke1UhwkP6URJ8JODgBcRRGpMaPFxiH/mEA6sxuu/+TvFf6ZRcnb6twBd3yRU2hDhNQk=
}
