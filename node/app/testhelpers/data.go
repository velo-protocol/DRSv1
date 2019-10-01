package testhelpers

import (
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

const (
	PublicKey1 = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
	SecretKey1 = "SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX"
	PublicKey2 = "GC2ROYZQH5FTVEPQZF7CAB32SCJC7DWVKILDUAT5BCU5O7HEI7HFUB25"
	SecretKey2 = "SCHQI345PYWHM2APNR4MN433HNCBS7VDUROOZKTYHZUBBTHI2YHNCJ4G"
	PublicKey3 = "GBGHQCINPG2257EN35E7EZA3D36KVXGSRNOVRCL6ERSWH2BIYQ5YUZKV"
	SecretKey3 = "SC6FARDSIVUTEYEYZ4KTO54RK5J5KRB2EW4JQ7QA5UFD3ZMYALTCYN5Y"

	DrsPublicKey                = "GCQCXIDTFMIL4VOAXWUQNRAMC46TTJDHZ3DDJVD32ND7B4OKANIUKB5N"
	DrsSecretKey                = "SDE374OE44ZU73KAUFYPNMQEUGCDIJLTIIUZ4W2MKWBPPAK36ID26ECU"
	TrustedPartnerListPublicKey = "GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"
	TrustedPartnerListSecretKey = "SBA7D4DPZYVRYMO7CO5EW57YEQNTUUI5DZXLULDVOX6DAYMRHF4JRVMH"
	RegulatorListPublicKey      = "GBEZLZLQKGH7I3VBU4VHWCAKFY6PPEN4DF7UDQO2AI7RELMKPLIOOODW"
	RegulatorListSecretKey      = "SBUAFRBY2CUPQKVFHGUAGODDKCKRXTL5RKQBPPGDABNIHDYJKJ7X4BZW"
	PriceFeederListPublicKey    = "GAP425PLI2M6ZPPC4A2D5DX45DR4AODPKLZWXPOO2XTICGTKGWLDGULL"
	PriceFeederListSecretKey    = "SAGIOUFA4DNRUYD64NGOAG4ETWNUEKURUYUB35IGHHAAXVLCXEPDN73T"
)

var (
	Kp1, _ = vconvert.SecretKeyToKeyPair(SecretKey1)
	Kp2, _ = vconvert.SecretKeyToKeyPair(SecretKey2)
	Kp3, _ = vconvert.SecretKeyToKeyPair(SecretKey3)

	DrsKp, _                = vconvert.SecretKeyToKeyPair(DrsSecretKey)
	TrustedPartnerListKp, _ = vconvert.SecretKeyToKeyPair(TrustedPartnerListSecretKey)
	RegulatorListKp, _      = vconvert.SecretKeyToKeyPair(RegulatorListSecretKey)
	PriceFeederListKp, _    = vconvert.SecretKeyToKeyPair(PriceFeederListSecretKey)
)

var (
	DrsAccountDataEntity = entities.DrsAccountData{
		TrustedPartnerListAddress: TrustedPartnerListPublicKey,
		RegulatorListAddress:      RegulatorListPublicKey,
		PriceFeederListAddress:    PriceFeederListPublicKey,
		PriceUsdVeloAddress:       "",
		PriceThbVeloAddress:       "",
		PriceSgdVeloAddress:       "",
		Base64Decoded:             true,
	}
)
