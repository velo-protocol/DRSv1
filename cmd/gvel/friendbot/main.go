package friendbot

type fb struct {
	FriendbotURL string
}

func NewFriendbot(friendbotURL string) Repository {
	return &fb{
		FriendbotURL: friendbotURL,
	}
}
