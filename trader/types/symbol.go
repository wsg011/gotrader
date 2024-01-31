package types

type SymbolInfo struct {
	Base       string  // 基础货币
	Quote      string  // 计价货币
	Symbol     string  // 交易对标识符
	FaceVal    float64 // 合约面值
	Multiplier int64   // 合约倍数
	PxPrec     int32   // 价格精度
	QtyPrec    int32   // 数量精度
	MinCnt     float64 // 最小交易量
	MaxCnt     float64 // 最大交易量
	Name       string  // 交易对名称
}
