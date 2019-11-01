package parser

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"net/url"
)

func ParseHorizonError(err error, horizonUrl string, network string) error {
	if herr, ok := err.(*horizonclient.Error); ok {
		errMsg := herr.Problem.Detail

		if envelope, _ := herr.EnvelopeXDR(); envelope != "" {
			if horizonUrl == "" || network == "" {
				errMsg = errMsg + " " + envelope
			} else {
				errMsg = errMsg + " " + makeXdrViewerUrl(envelope, horizonUrl, network)
			}
		}

		return errors.New(errMsg)
	}

	return err
}

func makeXdrViewerUrl(envelope string, horizonUrl string, network string) string {
	query := (&url.URL{}).Query()
	query.Add("input", envelope)
	query.Add("type", "TransactionEnvelope")
	query.Add("network", "custom")
	query.Add("horizonURL", horizonUrl)
	query.Add("networkPassphrase", network)

	return (&url.URL{
		Scheme:   "https",
		Host:     "www.stellar.org",
		Path:     "laboratory",
		Fragment: "xdr-viewer?" + query.Encode(),
	}).String()
}
