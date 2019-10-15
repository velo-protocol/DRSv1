package friendbot

type Repository interface {
	GetFreeLumens(stellarAddress string) error
}
