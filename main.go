package main

import (

	// HTTP
	"net/http"

	// Utils
	"fmt"
	"os"
	"strconv"

	// Redis
	"gopkg.in/redis.v5"

	. "github.com/hexagon/gotris/server"
)

// Version
const (
	appMajor uint = 0
	appMinor uint = 9
	appPatch uint = 1

	appPreRelease = ""
)

var (
	redisClient *redis.Client

	redisAddr  string
	redisPass  string
	listenPort string
	assetsPath string
)

func main() {

	fmt.Println(fmt.Sprintf("Gotris %s starting...", version()))

	// Read configuration from environment
	redisAddr = os.Getenv("GOTRIS_REDIS_ADDR")
	redisPass = os.Getenv("GOTRIS_REDIS_PASS")
	listenPort = os.Getenv("GOTRIS_PORT")
	assetsPath = os.Getenv("GOTRIS_ASSETS")

	// Apply defaults
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	if listenPort == "" {
		listenPort = "8080"
	}

	// Convert port from string to integer
	listenPortInt, portConvErr := strconv.Atoi(listenPort)
	if portConvErr != nil {
		fmt.Println("Invalid port: %d", listenPort)
		return
	}

	// Connect to redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass, // no password set
		DB:       0,         // use default DB
	})
	_, redisErr := redisClient.Ping().Result()

	// Bail out if redis connection failed
	if redisErr != nil {
		fmt.Println("Redis connection error: ", redisErr)
		return
	} else {
		fmt.Println("Connected to redis")
	}

	// Serve url /static from fs ./static/
	fs := http.FileServer(http.Dir(assetsPath))

	// Get template handler
	templateHandler := NewTemplateHandler(assetsPath, version())
	websocketHandler := NewWSHandler(redisClient)

	// Handlers
	http.Handle("/static/", fs)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/api/highscores", HighscoreHandler(redisClient))
	http.HandleFunc("/ws", websocketHandler)
	http.HandleFunc("/", templateHandler)

	// Listen to tcp port
	fmt.Println(fmt.Sprintf("Listening on *:%d...", listenPortInt))
	err := http.ListenAndServe(fmt.Sprintf(":%d", listenPortInt), nil)
	if err != nil {
		fmt.Println("Fatal error:", err)
	}

}

func version() string {
	if appPreRelease != "" {
		return fmt.Sprintf("%d.%d.%d-%s", appMajor, appMinor, appPatch, appPreRelease)
	} else {
		return fmt.Sprintf("%d.%d.%d", appMajor, appMinor, appPatch)
	}
}
