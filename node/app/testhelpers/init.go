package testhelpers

import (
	"github.com/velo-protocol/DRSv1/node/app/environments"
	"time"
)

func InitEnv() {
	env.DrsPublicKey = "GCQCXIDTFMIL4VOAXWUQNRAMC46TTJDHZ3DDJVD32ND7B4OKANIUKB5N"
	env.DrsSecretKey = "SDE374OE44ZU73KAUFYPNMQEUGCDIJLTIIUZ4W2MKWBPPAK36ID26ECU"
	env.VeloIssuerPublicKey = "GCV3Q6QZZUG7RWPIQ5CZ5MOW3KBBCQYT64EIQX6GVKDQL27WOYNDQD3G"
	env.NetworkPassphrase = "Test SDF Network ; September 2015"
	env.HorizonURL = "https://horizon.com"
	env.ValidPriceBoundary = 15 * time.Minute
}
