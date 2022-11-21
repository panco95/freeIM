package utils

import "github.com/tkuchiki/go-timezone"

func GetTimezones() []string {
	arr := []string{}
	tz := timezone.New()
	info := tz.TzInfos()
	for k := range info {
		arr = append(arr, k)
	}
	return arr
}
