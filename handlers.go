package main

import (

	// HTTP
	"html/template"
	"net/http"

	// Util
	"fmt"
	"os"
	"path/filepath"

	// Redis
	"gopkg.in/redis.v5"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   8192,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	maxMessageSize int64 = 1024 * 10 // Websocket request cannot be larger than 10KB

)

func TemplateHandler(w http.ResponseWriter, r *http.Request) {

	// Default to [..]/index.html by checking if last character is /
	up := filepath.Clean(r.URL.Path)
	if up[len(up)-1:] == "/" {
		up = up + "index.html"
	}

	lp := filepath.Join(assetsPath, "templates", "layout.html")
	fp := filepath.Join(assetsPath, "templates", up)

	fmt.Println("%s", lp)

	// Return a 404 if the base template doesn't exist
	info, err := os.Stat(lp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// Return a 404 if the requested template doesn't exist
	info, err = os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// Compose base template (layout.html) and (index.html)
	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		fmt.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Render template with data from pd
	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

}

func WSHandler(w http.ResponseWriter, r *http.Request) {

	// Try to upgrade connection
	c, err := upgrader.Upgrade(w, r, nil)

	// Defer connection close before doing anything else
	defer func() {
		fmt.Print("Closing connection")
		c.Close()
	}()

	// Check if stuff went wrong
	if err != nil {
		fmt.Print("HTTP Upgrade Error:", err)
		return
	}

	// Start a Client worker
	// We do not need to spawn a new goroutine as this handler
	// already are in it's own goroutine
	Client(c)

}

func HighscoreHandler(redisClient *redis.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Send current highscore
		// ToDo: Move to a separate handler and get by ajax
		// Send high score
		vals, rvErr := redisClient.ZRevRangeByScoreWithScores("gotris", redis.ZRangeBy{
			Min:    "1000",
			Max:    "+inf",
			Offset: 0,
			Count:  10,
		}).Result()

		if rvErr != nil {

			// An error occurred, response code 500
			fmt.Println("Redis get error: ", rvErr)
			http.Error(w, http.StatusText(500), 500)

		} else {

			// Success, response code 200, json
			highScoreJson := "{ \"highscore\": ["
			first := true
			for _, hs := range vals {
				if !first {
					highScoreJson += ","
				} else {
					first = false
				}
				highScoreJson += hs.Member.(string)

			}
			highScoreJson += "] }"
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(highScoreJson))
		}
	}
}
