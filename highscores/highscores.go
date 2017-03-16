package highscores

import (
	// Util
	"encoding/json"
	"fmt"
	"time"

	// mgo MongoDB driver
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Highscore struct {
	Nickname string
	Score    uint64
	Level    uint64
	Lines    uint64
	Ts       time.Time
}

type HighscoreMessage struct {
	Ath  []Highscore
	Week []Highscore
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day-7, 0, 0, 0, 0, t.Location())
}

func Write(session *mgo.Session, h Highscore) bool {

	c := session.DB("gotris").C("highscore")

	err := c.Insert(&h)

	if err != nil {
		fmt.Println("Error writing highscores", err)
		return false
	} else {
		return true
	}

}

func Read(session *mgo.Session) ([]byte, error) {

	c := session.DB("gotris").C("highscore")

	var ath []Highscore
	var week []Highscore

	errAth := c.Find(bson.M{}).Sort("-score").Limit(9).All(&ath)
	errWeek := c.Find(bson.M{"ts": bson.M{"$gt": Bod(time.Now())}}).Sort("-score").Limit(9).All(&week)

	if errAth != nil {
		return nil, errAth
	} else if errWeek != nil {
		return nil, errWeek
	}

	return json.Marshal(HighscoreMessage{ath, week})

}
