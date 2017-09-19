package main

import (

	// HTTP
	"net/http"

	// Utils
	"fmt"
	"os"
	"strconv"
	"time"

	// Redis
	"gopkg.in/mgo.v2"

	. "github.com/hexagon/gotris/server"
)

// Version
const (
	appMajor uint = 0
	appMinor uint = 9
	appPatch uint = 5

	appPreRelease = ""
)

var (
	mongoAddr  string
	mongoUser  string
	mongoPass  string
	mongoDB    string
	listenPort string
	assetsPath string
)

func main() {

	fmt.Println(fmt.Sprintf("Gotris %s starting...", version()))

	// Read configuration from environment
	mongoAddr = os.Getenv("GOTRIS_MONGO_ADDR")
	mongoUser = os.Getenv("GOTRIS_MONGO_USER")
	mongoPass = os.Getenv("GOTRIS_MONGO_PASS")
	mongoDB = os.Getenv("GOTRIS_MONGO_DB")
	listenPort = os.Getenv("GOTRIS_PORT")
	assetsPath = os.Getenv("GOTRIS_ASSETS")

	// Apply defaults
	if mongoAddr == "" {
		mongoAddr = "127.0.0.1"
	}
	if listenPort == "" {
		listenPort = "8080"
	}
	if mongoDB == "" {
		mongoDB = "gotris"
	}

	// Convert port from string to integer
	listenPortInt, portConvErr := strconv.Atoi(listenPort)
	if portConvErr != nil {
		fmt.Println("Invalid port: %d", listenPort)
		return
	}

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoAddr},
		Timeout:  60 * time.Second,
		Database: mongoDB,
		Username: mongoUser,
		Password: mongoPass,
	}

	mgoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()

	mgoSession.SetMode(mgo.Monotonic, true)

	// Serve url /static from fs ./static/
	fs := http.FileServer(http.Dir(assetsPath))

	// Get template handler
	templateHandler := NewTemplateHandler(assetsPath, version())
	websocketHandler := NewWSHandler(mgoSession)

	// Handlers
	http.Handle("/static/", fs)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/api/highscores", HighscoreHandler(mgoSession))
	http.HandleFunc("/ws", websocketHandler)
	http.HandleFunc("/", templateHandler)

	// Listen to tcp port
	fmt.Println(fmt.Sprintf("Listening on *:%d...", listenPortInt))
	err = http.ListenAndServe(fmt.Sprintf(":%d", listenPortInt), nil)
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
