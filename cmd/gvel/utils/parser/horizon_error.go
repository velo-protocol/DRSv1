package parser

import (
	"fmt"
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

	fragmentQuery := fmt.Sprintf("?input=%s&type=%s&network=%s&horizonURL=%s&networkPassphrase=%s", url.QueryEscape(envelope), "TransactionEnvelope", "custom", horizonUrl, network)

	return (&url.URL{
		Scheme:   "https",
		Host:     "www.stellar.org",
		Path:     "laboratory",
		Fragment: "xdr-viewer",
	}).String() + fragmentQuery // Todo: cannot push it inside url because it repeat encode with .String() function, refactor if we can
}
