package constants

import "github.com/pkg/errors"

var (
	ErrCreateWhiteList = errors.New("can't create white list")
	ErrUnauthorized = errors.New("unauthorized to perform an action")
	ErrRoleNotFound = errors.New("role not found")
	ErrBadSignature = errors.New("bad signature")

	ErrToSaveDatabase      = errors.New("can't save to database")
	ErrToUpdateDatabase    = errors.New("can't update to database")
	ErrToBeginTransaction  = errors.New("can't start database transaction")
	ErrToCommitTransaction = errors.New("can't commit database transaction")
	ErrToDeleteRecord      = errors.New("can't delete record")
)
