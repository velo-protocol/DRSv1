package constants

import "github.com/pkg/errors"

var (
	ErrCreateWhiteList = errors.New("can't create white list")

	ErrRoleNotFound = errors.New("role not found")

	ErrorToSaveDatabase      = errors.New("can't save to database")
	ErrorToUpdateDatabase    = errors.New("can't update to database")
	ErrorToBeginTransaction  = errors.New("can't start database transaction")
	ErrorToCommitTransaction = errors.New("can't commit database transaction")
	ErrorToDeleteRecord      = errors.New("can't delete record")
)
