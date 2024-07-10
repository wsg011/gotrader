package marketcenter

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/wsg011/gotrader/trader/constant"
	"github.com/wsg011/gotrader/trader/types"
)

func handleTicker(data interface{}) {
	if ticker, ok := data.(types.BookTicker); ok {
		fmt.Printf("BookTicker Update - Symbol: %s, Price: %.2f\n", ticker.Symbol, ticker.AskPrice)
	}
}

func simulateTickers(center *MarketCenter) {
	symbols := []string{"BTC_USDT", "ETH_USDT", "SOL_USDT"}                              // 示例股票代码
	exchanges := []constant.ExchangeType{constant.BinanceSpot, constant.BinanceUFutures} // 示例交易所

	for {
		for _, exchange := range exchanges {
			for _, symbol := range symbols {
				ticker := types.BookTicker{
					Exchange: exchange,
					Symbol:   symbol,
					AskPrice: 100 + rand.Float64()*50, // 随机交易量
				}
				center.Publish(ticker)
			}
		}
		time.Sleep(time.Second) // 每秒更新一次数据
	}
}

func TestMarketCenter(t *testing.T) {
	center := NewMarketCenter()

	// 订阅AAPL的Ticker数据
	center.Subscribe(constant.BinanceSpot.Name(), "BTC_USDT", "BookTicker", handleTicker)
	center.Subscribe(constant.BinanceSpot.Name(), "ETH_USDT", "BookTicker", handleTicker)

	// 模拟数据生成和发布（这里应该在新的goroutine中运行）
	go simulateTickers(center)

	// 模拟数据函数需要调整，确保发布数据时指定Exchange, Symbol, 和 Topic
	select {} // 阻塞main函数，保持程序运行
}
