package highscores

import (
	// Util
	"encoding/json"
	"fmt"
	"time"

	// Redis
	"gopkg.in/redis.v5"
)

type Highscore struct {
	Nickname string
	Score    uint64
	Level    uint64
	Lines    uint64
	Ts       time.Time
}

func (h *Highscore) Marshal() ([]byte, error) {
	return json.Marshal(&h)
}

func Write(client *redis.Client, h Highscore) bool {

	hMarshaled, marshalErr := h.Marshal()

	if marshalErr != nil {
		fmt.Println("Marshal error", marshalErr)
		return false
	}

	entry := redis.Z{
		Score:  float64(h.Score),
		Member: hMarshaled,
	}

	err := client.ZAdd("gotris", entry).Err()
	if err != nil {
		fmt.Println("Error writing highscores", err)
		return false
	} else {
		return true
	}

}
