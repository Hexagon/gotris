package server

import (

	// Util
	"fmt"
	"time"

	//
	"github.com/hexagon/gotris/game"
	"github.com/hexagon/gotris/highscores"

	// Redis
	"gopkg.in/redis.v5"
)

func Player(wsInChannel chan map[string]interface{}, wsOutChannel chan string, redisClient *redis.Client) {

	defer func() {
		fmt.Println("Exiting Player")
		close(wsOutChannel)
	}()

	// Set up game instance
	var g game.Game

	// Player is not ready yet!
	ready := false

	for {

		// Incoming messages
		select {
		case packet, packetOk := <-wsInChannel:
			if packetOk {
				switch packet["packet"] {

				case "ready":
					if n, ok := packet["nickname"]; ok {
						g = game.NewGame(n.(string), wsOutChannel)
						ready = true
					} else {
						fmt.Println("Could not start game, nickname was missing")
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
			// If all channels are empty, wait for a slight delay
			time.Sleep(2 * time.Millisecond)
		}

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

	}

	fmt.Println("Game over for", g.Nickname, "at", g.Score)

}
