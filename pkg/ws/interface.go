package ws

type WsImp interface {
	Subscribe(symbol string, topic string)
	OnConnected()
	MsgHandle([]byte)
}
