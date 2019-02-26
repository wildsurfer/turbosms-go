package turbosms_go

import (
	"turbosms_go/wsdl"
)

type Client struct {
	SoapService *wsdl.ServiceSoap
	Debug       bool
}

func (c *Client) GetCreditBalance() (*wsdl.GetCreditBalanceResponse, error) {
	wsdl.Debug = c.Debug
	return c.SoapService.GetCreditBalance(&wsdl.GetCreditBalance{})
}

func (c *Client) SendSMS(sender string, destination string, text string, wappush string) (*wsdl.SendSMSResponse, error) {
	wsdl.Debug = c.Debug
	msg := &wsdl.SendSMS{
		Sender:      sender,
		Destination: destination,
		Text:        text,
		Wappush:     wappush,
	}
	return c.SoapService.SendSMS(msg)
}

func (c *Client) GetNewMessages() (*wsdl.GetNewMessagesResponse, error) {
	wsdl.Debug = c.Debug
	return c.SoapService.GetNewMessages(&wsdl.GetNewMessages{})
}

func (c *Client) GetMessageStatus(msgId string) (*wsdl.GetMessageStatusResponse, error) {
	wsdl.Debug = c.Debug
	return c.SoapService.GetMessageStatus(&wsdl.GetMessageStatus{MessageId: msgId})
}

func NewClient(username string, password string) *Client {
	soapService := wsdl.NewServiceSoap("", false, &wsdl.BasicAuth{})
	soapService.Auth(&wsdl.Auth{Login: username, Password: password})
	cli := Client{soapService, false}
	return &cli
}
