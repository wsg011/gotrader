package types

type Position struct {
	MarginMode       string  // 保证金模式
	Symbol           string  // 交易对或合约标识符
	LiquidationPx    float64 // 清算价格
	Side             string  // 持仓方向，LONG/SHOTY
	Position         float64 // 当前持仓量
	FrozenPosition   float64 // 冻结的持仓量
	AvgCost          float64 // 平均成本
	UnrealisedPnl    float64 // 未实现盈亏
	Last             float64 // 最后成交价格
	MaintMarginRatio float64 // 维持保证金率
	Margin           float64 // 已使用的保证金
	Leverage         float64 // 杠杆倍数
}
