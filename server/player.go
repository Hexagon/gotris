package server

import (

	// Util
	"fmt"
	"runtime"
	"time"

	//
	"github.com/hexagon/gotris/game"
	"github.com/hexagon/gotris/highscores"

	// Redis
	"gopkg.in/redis.v5"
)

func validateNickname(n string) (string, string) {
	if len(n) < 2 {
		return n, "Nickname too short, need to be at least 2 characters"
	} else if len(n) > 15 {
		return n, "Nickname too long, need to be at most 15 characters"
	} else {
		return n, ""
	}
}

func Player(wsInChannel chan map[string]interface{}, wsOutChannel chan string, redisClient *redis.Client) {

	defer func() {
		fmt.Println("Exiting Player")
		close(wsOutChannel)
	}()

	// Set up game instance
	var g game.Game

	// Player is not ready yet!
	ready := false

	// Incoming messages
	for {
		select {
		case packet, packetOk := <-wsInChannel:
			if packetOk {
				switch packet["packet"] {

				case "ready":
					if n, ok := packet["nickname"]; ok {
						nickname, err := validateNickname(n.(string))
						if err != "" {
							wsOutChannel <- "{ \"ready\": false, \"error\": \"" + err + "\" }"
						} else {
							g = game.NewGame(nickname, wsOutChannel)
							ready = true
							wsOutChannel <- "{ \"ready\": true }"
						}
					} else {
						fmt.Println("Could not start game, nickname was missing")
						wsOutChannel <- "{ \"ready\": false, \"error\": \"Missing or invalid nickname\" }"
					}

				case "key":
					if d, ok := packet["data"]; ok && ready {
						g.SetKey(d.(map[string]interface{}))
					} else {
						fmt.Println("Received key update was invalid")
					}
				}
			} else {
				return
			}
		default:

			if ready {

				// Always iterate gamefield
				if !g.Iterate() {

					// End condition
					wsOutChannel <- "{ \"gameOver\": true }"

					// Write highscore
					highscores.Write(redisClient, highscores.Highscore{
						Nickname: g.Nickname,
						Score:    g.Score,
						Level:    g.Level,
						Lines:    g.Lines,
						Ts:       time.Now(),
					})

					break

				}

			}

			runtime.Gosched()
			time.Sleep(1500 * time.Microsecond)

		}

	}

	fmt.Println("Game over for", g.Nickname, "at", g.Score)

}
