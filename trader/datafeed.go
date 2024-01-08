package trader

type DataFeed struct {
	PublicDataChan  chan MarketDataInterface // 公开市场数据通道
	PrivateDataChan chan MarketDataInterface // 私有市场数据通道
}

func NewDataFeed() *DataFeed {
	return &DataFeed{
		PublicDataChan:  make(chan MarketDataInterface, 100),
		PrivateDataChan: make(chan MarketDataInterface, 100),
	}
}

// 接收数据并根据类型分发到相应通道的方法
func (feed *DataFeed) ReceiveData(data MarketDataInterface) {
	switch data.(type) {
	case *BookTicker, *OrderBook:
		feed.PublicDataChan <- data
	case *Trade:
		feed.PrivateDataChan <- data
	}
}
