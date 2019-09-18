package constants

import "github.com/pkg/errors"

var (
	ErrCreateWhiteList = errors.New("can't create white list")

	ErrRoleNotFound = errors.New("role not found")
	ErrRoleIsNotValid = errors.New("role is not valid")

	ErrToSaveDatabase      = errors.New("can't save to database")
	ErrToUpdateDatabase    = errors.New("can't update to database")
	ErrToBeginTransaction  = errors.New("can't start database transaction")
	ErrToCommitTransaction = errors.New("can't commit database transaction")
	ErrToDeleteRecord      = errors.New("can't delete record")
)
