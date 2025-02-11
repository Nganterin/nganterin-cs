package helpers

import "time"

func TimeToHumanReadable(t time.Time) string {
	return t.Format("15:04 Jan 2 2006")
}
