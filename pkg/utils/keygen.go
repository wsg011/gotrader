package utils

import (
	"math/rand"
	"sync"
	"time"
)

var (
	charset = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	keyMutex    sync.Mutex
	keyRander   *rand.Rand
	randerCache []byte
)

func init() {
	keyRander = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString 随机
func RandomString(n int) string {
	keyMutex.Lock()
	defer keyMutex.Unlock()
	if len(randerCache) < n {
		randerCache = make([]byte, n)
	}
	for index := 0; index < n; index++ {
		randerCache[index] = charset[keyRander.Intn(len(charset))]
	}
	return string(randerCache[:n])
}

// RandomString64 生成长度位64得随机字符串
func RandomString64() string {
	return RandomString(64)
}

// RandomStringer 不加锁,使用者自己保证线程安全
type RandomStringer struct {
	rander *rand.Rand
	buff   []byte
}

func (r *RandomStringer) RandomTracingID() string {
	return r.RandomString(10)
}

// RandomString 随机
func (r *RandomStringer) RandomString(n int) string {
	keyMutex.Lock()
	defer keyMutex.Unlock()
	if len(r.buff) < n {
		r.buff = make([]byte, n)
	}
	for index := 0; index < n; index++ {
		r.buff[index] = charset[r.rander.Intn(len(charset))]
	}
	return string(r.buff[:n])
}

func NewRandomStringer() *RandomStringer {
	return &RandomStringer{
		rander: rand.New(rand.NewSource(time.Now().UnixNano())),
		buff:   make([]byte, 16),
	}
}
