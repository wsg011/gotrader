package event

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestEventEngine(t *testing.T) {
	engine := NewEventEngine()
	engine.Register("testEvent", func(e Event) {
		fmt.Println("Handler received:", e.Data)
	})

	engine.Push("testEvent", "Hello, World!")

	// 等待事件处理完成
	time.Sleep(time.Second)
}

func TestEventEngineConcurrency(t *testing.T) {
	engine := NewEventEngine()

	var wg sync.WaitGroup
	// 创建一个channel用于接收事件数据
	dataChannel := make(chan int, 100)

	handler := func(e Event) {
		// 模拟一些处理时间
		time.Sleep(10 * time.Millisecond)
		wg.Done()
	}

	// 注册处理器
	engine.Register("concurrentEvent", handler)

	eventCount := 10000 // 测试的事件数量
	wg.Add(eventCount)

	start := time.Now()

	// 并发触发事件
	for i := 0; i < eventCount; i++ {
		engine.Push("concurrentEvent", i)
	}

	// 等待所有事件处理完成
	wg.Wait()
	close(dataChannel)
	duration := time.Since(start)

	t.Logf("Processed %d events in %v", eventCount, duration)
}
