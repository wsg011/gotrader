package types

type ExchangeParameters struct {
	DebugMode  bool
	ProxyURL   string // example: socks5://127.0.0.1:1080 | http://127.0.0.1:1080
	AccessKey  string
	SecretKey  string
	Passphrase string
}
