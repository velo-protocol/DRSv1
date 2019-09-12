package node

import "gitlab.com/velo-labs/cen/app/entities"

type Repository interface {
	SaveCredit(credit entities.Credit) error
}
