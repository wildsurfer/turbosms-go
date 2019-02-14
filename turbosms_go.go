package turbosms_go

import (
	"github.com/wildsurfer/turbosms-go/wsdl"
)

type Client struct {
	SoapService *wsdl.ServiceSoap
}

func (c *Client) GetCreditBalance() (*wsdl.GetCreditBalanceResponse, error) {
	return c.SoapService.GetCreditBalance(&wsdl.GetCreditBalance{})
}

func (c *Client) SendSMS(sender string, destination string, text string, wappush string) (*wsdl.SendSMSResponse, error) {
	msg := &wsdl.SendSMS{
		Sender:      sender,
		Destination: destination,
		Text:        text,
		Wappush:     wappush,
	}
	return c.SoapService.SendSMS(msg)
}

func (c *Client) GetNewMessages() (*wsdl.GetNewMessagesResponse, error) {
	return c.SoapService.GetNewMessages(&wsdl.GetNewMessages{})
}

func (c *Client) GetMessageStatus(msgId string) (*wsdl.GetMessageStatusResponse, error) {
	return c.SoapService.GetMessageStatus(&wsdl.GetMessageStatus{MessageId: msgId})
}

func NewClient(username string, password string) *Client {
	soapService := wsdl.NewServiceSoap("", false, &wsdl.BasicAuth{})
	soapService.Auth(&wsdl.Auth{Login: username, Password: password})
	cli := Client{soapService}
	return &cli
}
