package stellar

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

func (s *stellar) GetFreeLumens(stellarAddress string) error {
	resp, err := http.Get(fmt.Sprintf(s.FriendBotURL, stellarAddress))
	if err != nil {
		return errors.Wrap(err, "failed to get free lumens from friendbot")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to get free lumens from friendbot")
	}

	return nil
}
