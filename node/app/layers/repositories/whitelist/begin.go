package whitelist

import "github.com/jinzhu/gorm"

func (r *repo) BeginTx() *gorm.DB {
	return r.Conn.Begin()
}
