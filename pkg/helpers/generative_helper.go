package helpers

import "time"

func GenerateUniquePassword() string {
	return time.Now().Format("20060102150405") + GenerateMilliseconds()
}

func GenerateMilliseconds() string {
	return time.Now().Format(".000")[1:]
}
