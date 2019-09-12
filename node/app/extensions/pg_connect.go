package extensions

import (
	"fmt"
	"github.com/golang-migrate/migrate"
	"gitlab.com/velo-labs/cen/node/app/environments"

	// for postgres db
	_ "github.com/golang-migrate/migrate/database/postgres"
	// for postgres db
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// for migration file
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"

	"math"
)

// ConnectDB connect to DB
func ConnectDB() *gorm.DB {
	connection := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.PostgresHost,
		env.PostgresPort,
		env.PostgresUser,
		env.PostgresPassword,
		env.PostgresDB,
	)

	db, err := gorm.Open("postgres", connection)

	if err != nil {
		panic(fmt.Sprintf("connect DB error : %+v", err.Error()))
	}

	return db
}

// DBMigration run migration files
func DBMigration() {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.PostgresUser,
		env.PostgresPassword,
		env.PostgresHost,
		env.PostgresPort,
		env.PostgresDB,
	)

	migrationEngine, err := migrate.New("file://../app/migrations", databaseURL)

	if err != nil {
		panic(fmt.Sprintf("DBMigration error : %+v", err.Error()))
	}

	_ = migrationEngine.Steps(math.MaxInt32)
}
