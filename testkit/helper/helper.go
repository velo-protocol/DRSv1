package helper

import (
	"encoding/json"
	"fmt"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
)

const (
	PublicKeyFirstRegulator = "GAI2RZUCHITWE46ABRXZMUE23OTSPDBDYDUNZOMQNJAS5RMOZYW5ZOZS"
	SecretKeyFirstRegulator = "SDCTXCAIL2AVXZGVM5X6SHCUF7F7K2IOB6MHTPGDPD2632B34WJHVB5Z"

	PublicKeyTP  = "GATWZYWCUC6ZJZF2O745HBV5GE7DAVE6UOXX57XMFSY22EB67AUXA7SI"
	SecretKeyTP  = "SC5USFD5OPJL6GQEN2Q7QOUA4RAGOTEQD2ONX7R2I34NBU576YR6VDQX"
	PublicKeyReg = "GCIN3UIQUO7DZ73NZ6IOY3Z45C35VY3PGS537JQZ2JJFC73CLVDDKRGJ"
	SecretKeyReg = "SACDLB5E7UJMWMMWOF2PPKRWEFK2FVOYVEEEB4SJBOL7RQUQ6WPZWMZS"
	PublicKeyPF  = "GDD7Q3GMQ4Q6NMX52JSPIWWO746HIIF4ES7PML6PI6ESC7DBARGTWFYL"
	SecretKeyPF  = "SACUMSJ64ZG6CFI7WMXN7MGT35LGEY5OZDA7QHZZB57LIVEIE7KH2GCH"
)

var (
	KPFirstRegulator, _ = vconvert.SecretKeyToKeyPair(SecretKeyFirstRegulator)
	KPTP, _             = vconvert.SecretKeyToKeyPair(SecretKeyTP)
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
