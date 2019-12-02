package models

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/strkey"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/utils"
	"strconv"
	"strings"
)

type Price struct {
	Source        string
	Value         int64
	UnixTimestamp int64
}

func NewPrice(priceSourceAddress string, encodedValue string) (*Price, error) {
	if !strkey.IsValidEd25519PublicKey(priceSourceAddress) {
		return nil, errors.Errorf(constants.ErrToDecodeData, priceSourceAddress)
	}

	value, err := utils.DecodeBase64(encodedValue)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(constants.ErrToDecodeData, priceSourceAddress))
	}

	parts := strings.Split(value, "_")
	if len(parts) != 2 {
		return nil, errors.Errorf("fail to parse price data from %s", priceSourceAddress)
	}

	priceTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("fail to parse price timestamp from %s", priceSourceAddress))
	}

	priceValue, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("fail to parse price value from %s", priceSourceAddress))
	}

	return &Price{
		Source:        priceSourceAddress,
		Value:         priceValue,
		UnixTimestamp: priceTimestamp,
	}, nil
}
