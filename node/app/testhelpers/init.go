package testhelpers

import (
	"gitlab.com/velo-labs/cen/node/app/environments"
	"time"
)

func InitEnv() {
	env.DrsPublicKey = "GCQCXIDTFMIL4VOAXWUQNRAMC46TTJDHZ3DDJVD32ND7B4OKANIUKB5N"
	env.DrsSecretKey = "SDE374OE44ZU73KAUFYPNMQEUGCDIJLTIIUZ4W2MKWBPPAK36ID26ECU"
	env.NetworkPassphrase = "Test SDF Network ; September 2015"
	env.HorizonURL = "https://horizon.com"
	env.ValidPriceBoundary = 15 * time.Minute
}
