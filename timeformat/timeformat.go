package timeformat

import (
	"CP_Discussion/log"
	"time"
)

const localLocation = "Asia/Taipei"

var local = buildLocal()

func buildLocal() *time.Location {
	l, err := time.LoadLocation(localLocation)
	if err != nil {
		log.Warning.Println(err)
		return time.Local
	}
	return l
}

func FormatTime(t *time.Time) {
	*t = t.In(local)
}
