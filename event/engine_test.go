package event

import (
	"github.com/wsg011/gotrader/trader/constant"
	"testing"
	"time"
)

func TestEventEngine(t *testing.T) {
	engine := NewEventEngine()
	engine.Start()

	engine.Register("testEvent", func(data interface{}) {
		t.Logf("Handler received: %s", data)
	})

	engine.Register("binance", func(data interface{}) {
		t.Logf("binance received: %s", data)
	})

	engine.Register("ok", func(data interface{}) {
		t.Logf("ok received: %s", data)
	})

	engine.Register("ok-bookticker", func(data interface{}) {
		t.Logf("ok-bookticker received: %s", data)
	})

	engine.Register("binance-bookticker", func(data interface{}) {
		t.Logf("binance-bookticker received: %s", data)
	})

	engine.Push("testEvent", "Hello, World!")
	engine.Push(constant.EVENT_BOOKTICKER, "binance-bookticker-btcusdt")
	engine.Push("binance-bookticker-btcusdt", "binance-bookticker-btcusdt")
	engine.Push("binance-bookticker-ethusdt", "binance-bookticker-ethusdt")
	engine.Push("ok-bookticker-ethusdt", "ok-bookticker-ethusdt")
	engine.Push("ok-orders-ethusdt", "ok-orders-ethusdt")
	engine.Push("ok-ticker-ethusdt", "ok-ticker-ethusdt")

	// 等待事件处理完成
	time.Sleep(time.Second)
}

// func TestEventEngineConcurrency(t *testing.T) {
// 	engine := NewEventEngine()

// 	var wg sync.WaitGroup
// 	// 创建一个channel用于接收事件数据
// 	dataChannel := make(chan int, 100)

// 	handler := func(e interface{}) {
// 		// 模拟一些处理时间
// 		time.Sleep(10 * time.Millisecond)
// 		wg.Done()
// 	}

// 	// 注册处理器
// 	engine.Register("concurrentEvent", handler)

// 	eventCount := 10000 // 测试的事件数量
// 	wg.Add(eventCount)

// 	start := time.Now()

// 	// 并发触发事件
// 	for i := 0; i < eventCount; i++ {
// 		engine.Push("concurrentEvent", i)
// 	}

// 	// 等待所有事件处理完成
// 	wg.Wait()
// 	close(dataChannel)
// 	duration := time.Since(start)

// 	t.Logf("Processed %d events in %v", eventCount, duration)
// }
