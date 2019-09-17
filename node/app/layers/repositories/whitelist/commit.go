package whitelist

import (
	"github.com/jinzhu/gorm"
	verrors "gitlab.com/velo-labs/cen/libs/errors"
)

func (r *repo) CommitTx(dbtx *gorm.DB) error {
	if err := dbtx.Commit().Error; err != nil {
		return verrors.InternalError{Message: err.Error()}
	}

	return nil
}
