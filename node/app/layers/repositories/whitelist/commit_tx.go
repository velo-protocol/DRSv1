package whitelist

import (
	"github.com/jinzhu/gorm"
)

func (r *repo) CommitTx(dbTx *gorm.DB) error {
	if err := dbTx.Commit().Error; err != nil {
		return err
	}

	return nil
}
