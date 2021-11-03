package ali_mns

import "encoding/xml"

type AccountManager struct {
	cli     MNSClient
	decoder MNSDecoder
}

type OpenService struct {
	BaseResponse
	XMLName xml.Name `xml:"OpenService" json:"-"`
	OrderId string   `xml:"OrderId" json:"order_id"`
}

func NewAccountManager(client MNSClient) *AccountManager {
	return &AccountManager{
		cli:     client,
		decoder: NewAliMNSDecoder(),
	}
}

func (p *AccountManager) OpenService() (attr OpenService, err error) {
	_, err = send(p.cli, p.decoder, POST, nil, nil, "commonbuy/openservice", &attr)
	return
}
