package helper

import (
	"encoding/json"
	"fmt"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
)

const (
	PublicKeyFirstRegulator = "GAPERKAIHG6K6VPUS5ZVGV7KXEEA6WEC4SOFPS4I2HNUREI72PKFWBCC"
	SecretKeyFirstRegulator = "SBNRFCGZBCBIDFYWVM3NFNMY6E2C23VDPESHOX7RGSVQKXWE2T77JRW5"

	PublicKeyTrustedPartner = "GC2HMLPY4PUB4FSBIKXDXEGFW7XUZJ4UIR3MVEQET5YY4W34EH7SSZQF"
	SecretKeyTrustedPartner = "SCQQQP46ETSFWALXOP4FXXOYQAGK67YGXZUEOYOUKW3GIO2H7QRUIQGL"
	PublicKeyReg            = "GCIN3UIQUO7DZ73NZ6IOY3Z45C35VY3PGS537JQZ2JJFC73CLVDDKRGJ"
	SecretKeyReg            = "SACDLB5E7UJMWMMWOF2PPKRWEFK2FVOYVEEEB4SJBOL7RQUQ6WPZWMZS"
	PublicKeyPF             = "GDD7Q3GMQ4Q6NMX52JSPIWWO746HIIF4ES7PML6PI6ESC7DBARGTWFYL"
	SecretKeyPF             = "SACUMSJ64ZG6CFI7WMXN7MGT35LGEY5OZDA7QHZZB57LIVEIE7KH2GCH"
)

var (
	KPFirstRegulator, _ = vconvert.SecretKeyToKeyPair(SecretKeyFirstRegulator)
	KPTP, _             = vconvert.SecretKeyToKeyPair(SecretKeyTrustedPartner)
	KPPF, _             = vconvert.SecretKeyToKeyPair(SecretKeyPF)
)

func DecodeB64VeloTx(base64VeloTx string) {
	fmt.Println("##### Start Decode Base64 Velo Transaction #####")
	veloTx, err := vtxnbuild.TransactionFromXDR(base64VeloTx)
	if err != nil {
		panic(err)
	}
	veloTxByte, err := json.Marshal(veloTx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Velo Transaction: %s \n", string(veloTxByte))

	fmt.Println("##### End Decode Base64 Velo Transaction #####")
}

func CompareVeloTxSigner(base64VeloTx string, accountToCompare string) {
	fmt.Println("##### Start Compare Velo Transaction Signer #####")
	veloTx, err := vtxnbuild.TransactionFromXDR(base64VeloTx)
	if err != nil {
		panic(err)
	}
	compareAccount, err := vconvert.PublicKeyToKeyPair(accountToCompare)
	if err != nil {
		panic(err)
	}

	isMatch := veloTx.TxEnvelope().Signatures[0].Hint == compareAccount.Hint()

	fmt.Printf("Velo Tx Signer Account == %s : %+v \n", compareAccount.Address(), isMatch)

	fmt.Println("##### End Compare Velo Transaction Signer #####")
}
