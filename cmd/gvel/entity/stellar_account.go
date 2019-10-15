package entity

type StellarAccount struct {
	EncryptedSeed []byte `json:"encryptedSeed"`
	Nonce         []byte `json:"nonce"`
}
