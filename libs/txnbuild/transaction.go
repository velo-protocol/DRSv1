/*
Package vtxnbuild implements transactions and operations on the Velo Node.
This library provides an interface to the Velo transaction model.
For more information and further examples, see https://docs.velo.org/.
*/
package vtxnbuild

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/hash"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"github.com/velo-protocol/DRSv1/libs/xdr"
)

// VeloTx represents a Velo transaction.
type VeloTx struct {
	SourceAccount txnbuild.Account
	VeloOp        VeloOp

	veloXdrTx       vxdr.VeloTx
	veloXdrEnvelope *vxdr.VeloTxEnvelope
}

// Build for Transaction completely configures the Transaction. After calling Build,
// the Transaction is ready to be serialised or signed.
func (veloTx *VeloTx) Build() error {
	if veloTx.veloXdrEnvelope != nil {
		if veloTx.veloXdrEnvelope.Signatures != nil {
			return errors.New("transaction has already been signed, so cannot be rebuilt.")
		}
		// clear the existing XDR so we don't append to any existing fields
		veloTx.veloXdrEnvelope = &vxdr.VeloTxEnvelope{}
		veloTx.veloXdrEnvelope.VeloTx = vxdr.VeloTx{}
	}

	// reset veloXdrOp
	veloTx.veloXdrTx = vxdr.VeloTx{}

	if veloTx.SourceAccount == nil {
		return errors.New("sourceAccount must not be blank")
	}
	accountID := veloTx.SourceAccount.GetAccountID()
	if accountID == "" {
		return errors.New("sourceAccount must not be blank")
	}
	if !strkey.IsValidEd25519PublicKey(accountID) {
		return errors.New("invalid sourceAccount format")
	}

	_, err := keypair.Parse(accountID)
	if err != nil {
		return err
	}
	// Set account ID in XDR
	_ = veloTx.veloXdrTx.SourceAccount.SetAddress(accountID)

	// Assign VeloOp
	if err := veloTx.VeloOp.Validate(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("validation failed for %T operation", veloTx.VeloOp))
	}

	veloXdrOp, err := veloTx.VeloOp.BuildXDR()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to build operation %T", veloTx.VeloOp))
	}

	veloTx.veloXdrTx.VeloOp = veloXdrOp

	// Initialise transaction envelope
	veloTx.veloXdrEnvelope = &vxdr.VeloTxEnvelope{}
	veloTx.veloXdrEnvelope.VeloTx = veloTx.veloXdrTx

	return nil
}

// Sign for Transaction signs a previously built transaction. A signed transaction may be
// submitted to the Velo network.
func (veloTx *VeloTx) Sign(kps ...*keypair.Full) error {
	// Hash the transaction
	hashedValue, err := veloTx.Hash()
	if err != nil {
		return errors.Wrap(err, "failed to hashedValue transaction")
	}

	// Sign the hashedValue
	for _, kp := range kps {
		sig, err := kp.SignDecorated(hashedValue[:])
		if err != nil {
			return errors.Wrap(err, "failed to sign transaction")
		}
		// Append the signature to the envelope
		veloTx.veloXdrEnvelope.Signatures = append(veloTx.veloXdrEnvelope.Signatures, sig)
	}

	return nil
}

// Hash provides a signable object representing the Transaction on the specified Velo network.
func (veloTx *VeloTx) Hash() ([32]byte, error) {
	var txBytes bytes.Buffer

	_, err := xdr.Marshal(&txBytes, veloTx)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "marshal velo tx failed")
	}

	return hash.Hash(txBytes.Bytes()), nil
}

// Base64 returns the base 64 XDR representation of the transaction envelope.
func (veloTx *VeloTx) Base64() (string, error) {
	bs, err := veloTx.MarshalBinary()
	if err != nil {
		return "", errors.Wrap(err, "failed to get XDR bytestring")
	}

	return base64.StdEncoding.EncodeToString(bs), nil
}

// MarshalBinary returns the binary XDR representation of the transaction envelope.
func (veloTx *VeloTx) MarshalBinary() ([]byte, error) {
	var txBytes bytes.Buffer
	_, err := xdr.Marshal(&txBytes, veloTx.veloXdrEnvelope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal XDR")
	}

	return txBytes.Bytes(), nil
}

// BuildSignEncode performs all the steps to produce a final transaction suitable
// for submitting to the Velo network.
func (veloTx *VeloTx) BuildSignEncode(keyPairs ...*keypair.Full) (string, error) {
	err := veloTx.Build()
	if err != nil {
		return "", errors.Wrap(err, "couldn't build transaction")
	}

	err = veloTx.Sign(keyPairs...)
	if err != nil {
		return "", errors.Wrap(err, "couldn't sign transaction")
	}

	txeBase64, err := veloTx.Base64()
	if err != nil {
		return "", errors.Wrap(err, "couldn't encode transaction")
	}

	return txeBase64, err
}

// TransactionFromXDR parses the supplied transaction envelope in base64 XDR and returns a Transaction object.
func TransactionFromXDR(veloTxBase64 string) (VeloTx, error) {
	var veloXdrEnvelope vxdr.VeloTxEnvelope
	err := xdr.SafeUnmarshalBase64(veloTxBase64, &veloXdrEnvelope)
	if err != nil {
		return VeloTx{}, errors.Wrap(err, "the XDR message cannot be decoded")
	}

	var veloTx VeloTx
	veloTx.veloXdrTx = veloXdrEnvelope.VeloTx
	veloTx.veloXdrEnvelope = &veloXdrEnvelope

	veloTx.SourceAccount = &txnbuild.SimpleAccount{
		AccountID: veloXdrEnvelope.VeloTx.SourceAccount.Address(),
	}

	veloTx.VeloOp, err = operationFromXDR(veloXdrEnvelope.VeloTx.VeloOp)
	if err != nil {
		return VeloTx{}, err
	}

	return veloTx, nil
}

// TxEnvelope returns the TransactionEnvelope XDR struct.
func (veloTx *VeloTx) TxEnvelope() *vxdr.VeloTxEnvelope {
	return veloTx.veloXdrEnvelope
}
