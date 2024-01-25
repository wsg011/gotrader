package utils

import (
	"fmt"
	"math"
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

func FormatFloat(number float64, precision int) string {
	format := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(format, number)
}

// ScaleFloat 将float64数值缩放到指定的精度，并转换为字符串。
// number 是要缩放的数值，precision 是小数点后的位数，roundUp 默认向下取整up/down。
func ScaleFloat(number float64, precision int, roundUp string) string {
	scalingFactor := math.Pow(10, float64(precision))

	var scaledNumber float64
	if roundUp == "up" {
		scaledNumber = math.Ceil(number*scalingFactor) / scalingFactor
	} else {
		scaledNumber = math.Floor(number*scalingFactor) / scalingFactor
	}

	format := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(format, scaledNumber)
}
