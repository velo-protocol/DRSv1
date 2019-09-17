package role

import (
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type Repo interface {
	FindOne(role string) (*entities.Role, error)
}