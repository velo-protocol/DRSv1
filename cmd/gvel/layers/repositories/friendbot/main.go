package friendbot

type friendBot struct {
	FriendBotURL string
}

func NewFriendBot(friendBotURL string) Repository {
	return &friendBot{
		FriendBotURL: friendBotURL,
	}
}
