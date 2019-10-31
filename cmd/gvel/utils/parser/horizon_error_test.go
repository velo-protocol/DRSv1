package parser

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/support/render/problem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHorizonError(t *testing.T) {
	t.Run("success, typical error", func(t *testing.T) {
		err := errors.New("some error has occurred")
		parsedError := ParseHorizonError(err, "", "")
		assert.EqualError(t, err, parsedError.Error())
	})
	t.Run("success, horizon error without envelope", func(t *testing.T) {
		err := &horizonclient.Error{
			Problem: problem.NotFound,
		}
		parsedError := ParseHorizonError(err, "", "")
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
					"envelope_xdr": "AAAA...",
				},
			},
		}

		parsedError := ParseHorizonError(err, "https://fake-horizon.com", "Fake Net")
		assert.EqualError(t, parsedError, "Transaction fail "+makeXdrViewerUrl("AAAA...", "https://fake-horizon.com", "Fake Net"))
	})
}
