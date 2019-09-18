package testhelpers

import (
	env "gitlab.com/velo-labs/cen/node/app/environments"
	"os"
)

func InitEnv() {
	_ = os.Setenv("DRS_PUBLIC_KEY", "GD7R43KMK3AANO4TW722AKX6HZ7TKHKFZM5N4ASRUVU4FHB55V2JKOS2")
	_ = os.Setenv("DRS_SECRET_KEY", "SB3P4Z4XRC4TPPJ5JCVKFRFJ4GIOOQ34F5ZPKEI37WZSJUVQN66TQGHB")
	_ = os.Setenv("VELO_ISSUER_PUBLIC_KEY", "GABB2XA6ROY6IZWP2EKD3JIW3KMT6EYAWJMPR7DYWQ3F64BEC7NPLEQS")
	_ = os.Setenv("NETWORK_PASSPHRASE", "Test SDF Network ; September 2015")
	_ = os.Setenv("HORIZON_URL", "https://horizon.com")

	env.Init()
}