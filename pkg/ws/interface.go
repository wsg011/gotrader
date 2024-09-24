package ws

type ConnectType int

const (
	Connect ConnectType = iota
	Reconnect
)

type WsImp interface {
	Subscribe(symbol string, topic string) map[string]interface{}
	OnConnected(*WsClient, ConnectType)
	Handle(*WsClient, []byte)
	Ping(*WsClient)
}
