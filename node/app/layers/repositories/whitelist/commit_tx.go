package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/constants"
)

func (r *repo) CommitTx(dbTx *gorm.DB) error {
	if err := dbTx.Commit().Error; err != nil {
		return constants.ErrToCommitTransaction
	}

	return nil
}
