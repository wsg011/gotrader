package ws

type ConnectType int

const (
	Connect ConnectType = iota
	Reconnect
)

type WsImp interface {
	OnConnected(*WsClient, ConnectType)
	Handle(*WsClient, []byte)
	Ping(*WsClient)
}
