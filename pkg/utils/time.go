package utils

import "time"

func Millisec(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func Microsec(t time.Time) int64 {
	return t.UnixNano() / int64(time.Microsecond)
}

// IsoTime eg: 2018-03-16T18:02:48.284Z
func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}
