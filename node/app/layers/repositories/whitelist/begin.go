package whitelist

import "github.com/jinzhu/gorm"

func (r *repo) BeginTx() (*gorm.DB, error) {
	tx := r.Conn.Begin()
	return tx, nil
}
