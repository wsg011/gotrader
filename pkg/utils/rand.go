package utils

import (
	mrand "math/rand"
	"sync"
	"time"
)

var randerMutex sync.Mutex
var rander *mrand.Rand

func init() {
	rander = mrand.New(mrand.NewSource(time.Now().UnixNano()))
}

func GenerateRangeNum(min, max int64) int64 {
	randerMutex.Lock()
	defer randerMutex.Unlock()
	randNum := rander.Int63n(max-min) + min
	return randNum
}
