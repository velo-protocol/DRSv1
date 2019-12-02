package parser

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/support/render/problem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHorizonError(t *testing.T) {
	t.Run("success, typical error", func(t *testing.T) {
		err := errors.New("some error has occurred")
		parsedError := ParseHorizonError(err, "", "")
		assert.NotNil(t, parsedError)
		assert.EqualError(t, err, parsedError.Error())
	})
	t.Run("success, horizon error without envelope", func(t *testing.T) {
		err := &horizonclient.Error{
			Problem: problem.NotFound,
		}
		parsedError := ParseHorizonError(err, "", "")
		assert.NotNil(t, parsedError)
		assert.Contains(t, parsedError.Error(), problem.NotFound.Detail)
	})
	t.Run("success, horizon error with envelope but no horizon url and network provided", func(t *testing.T) {
		err := &horizonclient.Error{
			Problem: problem.P{
				Detail: "Transaction fail",
				Extras: map[string]interface{}{
					"envelope_xdr": "AAAA...",
				},
			},
		}

		parsedError := ParseHorizonError(err, "", "")
		assert.EqualError(t, parsedError, "Transaction fail AAAA...")
	})
	t.Run("success, horizon error with envelope, horizon url and network provided", func(t *testing.T) {
		err := &horizonclient.Error{
			Problem: problem.P{
				Detail: "Transaction fail",
				Extras: map[string]interface{}{
					"envelope_xdr": "AAAAALVxVHrmig1kCoVWZXIbF+JroljcpDNst+8lCC2UwFbrAAAAZAAAAjMAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAtHYt+OPoHhZBQq47kMW370ynlER2ypIEn3GOW3wh/ykAAAABVkVMTwAAAAC2JwYElOJJ6BdXQDj7a/rcwmbDe4duF1n+OcnMjpxDJQAAAAA7msoAAAAAAAAAAAGUwFbrAAAAQNHC2BigiO5Z5/RyngbiCYiyNZelgjnYpFtLqRd7wwMqNAWzyi3exojOus0nB28c3B87yT7qf9NP/AxHICqFVA4=",
				},
			},
		}

		parsedError := ParseHorizonError(err, "https://horizon-testnet.stellar.org", network.TestNetworkPassphrase)
		assert.Error(t, parsedError)
		assert.Equal(t, "Transaction fail https://www.stellar.org/laboratory#xdr-viewer?input=AAAAALVxVHrmig1kCoVWZXIbF%2BJroljcpDNst%2B8lCC2UwFbrAAAAZAAAAjMAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAtHYt%2BOPoHhZBQq47kMW370ynlER2ypIEn3GOW3wh%2FykAAAABVkVMTwAAAAC2JwYElOJJ6BdXQDj7a%2FrcwmbDe4duF1n%2BOcnMjpxDJQAAAAA7msoAAAAAAAAAAAGUwFbrAAAAQNHC2BigiO5Z5%2FRyngbiCYiyNZelgjnYpFtLqRd7wwMqNAWzyi3exojOus0nB28c3B87yT7qf9NP%2FAxHICqFVA4%3D&type=TransactionEnvelope&network=custom&horizonURL=https://horizon-testnet.stellar.org&networkPassphrase=Test SDF Network ; September 2015", parsedError.Error())
	})
}
