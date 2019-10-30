package console

import (
	"os"
)

const (
	// ExitSuccess code represents exit when success
	ExitSuccess = iota
	// ExitError code represents exit when error occur
	ExitError
	// ExitBadConnection code represents exit when bad connection
	ExitBadConnection
	// ExitInvalidInput code represents for txn, watch command
	ExitInvalidInput
	// ExitBadFeature code represents provided a valid flag with an unsupported value
	ExitBadFeature
	// ExitInterrupted code represents exit the code get interrupted
	ExitInterrupted
	// ExitIO code represents exit when io error
	ExitIO
	// ExitBadArgs code represents  exit when arguments are not right
	ExitBadArgs = 128
)

// ExitWithError Exit Command with error
func ExitWithError(code int, err error) {
	StopLoading()
	Logger.Error(err.Error())
	os.Exit(code)
}
