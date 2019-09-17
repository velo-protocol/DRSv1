package whitelist

import (
	"github.com/jinzhu/gorm"
	verrors "gitlab.com/velo-labs/cen/libs/errors"
)

func (r *repo) CommitTx(dbtx *gorm.DB) (bool, error) {
	if err := dbtx.Commit().Error; err != nil {
		return false, verrors.InternalError{Message: err.Error()}
	}

	return true, nil
}
