package okxv5swap

type OkxV5Swap struct {
	client *RestClient
}

func NewOkxV5Swap(apiKey, secretKey, passPhrase string) *OkxV5Swap {
	client := NewRestClient(apiKey, secretKey, passPhrase)

	return &OkxV5Swap{
		client: client,
	}

}

func (okx *OkxV5Swap) GetTickers() {

}
