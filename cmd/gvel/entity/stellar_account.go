package entity

type StellarAccount struct {
	Address       string `json:"address"`
	EncryptedSeed []byte `json:"encryptedSeed"`
	Nonce         []byte `json:"nonce"`
	IsDefault     bool   `json:"-"`
}
