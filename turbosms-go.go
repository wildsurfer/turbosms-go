// Package turbosms is a client for the legacy TurboSMS SOAP gateway.
//
// SOAP is the legacy TurboSMS protocol. For new projects consider the
// current HTTP API: https://turbosms.ua/en/api.html
package turbosms

import (
	"github.com/wildsurfer/turbosms-go/wsdl"
)

// Client is a TurboSMS SOAP gateway client. Create it with NewClient.
type Client struct {
	SoapService *wsdl.ServiceSoap
	// Debug enables logging of raw SOAP requests and responses.
	Debug bool
}

// GetCreditBalance returns the current account balance in credits.
func (c *Client) GetCreditBalance() (*wsdl.GetCreditBalanceResponse, error) {
	c.SoapService.SetDebug(c.Debug)
	return c.SoapService.GetCreditBalance(&wsdl.GetCreditBalance{})
}

// SendSMS sends a text message to the destination phone number.
//
// The gateway reports message level failures, for example not enough
// credits, as strings inside SendSMSResult. In this case the returned
// error is nil, so always check the result strings.
func (c *Client) SendSMS(sender string, destination string, text string, wappush string) (*wsdl.SendSMSResponse, error) {
	c.SoapService.SetDebug(c.Debug)
	msg := &wsdl.SendSMS{
		Sender:      sender,
		Destination: destination,
		Text:        text,
		Wappush:     wappush,
	}
	return c.SoapService.SendSMS(msg)
}

// GetNewMessages returns new incoming messages.
func (c *Client) GetNewMessages() (*wsdl.GetNewMessagesResponse, error) {
	c.SoapService.SetDebug(c.Debug)
	return c.SoapService.GetNewMessages(&wsdl.GetNewMessages{})
}

// GetMessageStatus returns the delivery status of the message with the given id.
func (c *Client) GetMessageStatus(msgId string) (*wsdl.GetMessageStatusResponse, error) {
	c.SoapService.SetDebug(c.Debug)
	return c.SoapService.GetMessageStatus(&wsdl.GetMessageStatus{MessageId: msgId})
}

// NewClient returns a Client that talks to the official TurboSMS SOAP
// endpoint using the given gateway credentials.
func NewClient(username string, password string) *Client {
	soapService := wsdl.NewServiceSoap("", false, &wsdl.BasicAuth{})
	soapService.Auth(&wsdl.Auth{Login: username, Password: password})
	cli := Client{soapService, false}
	return &cli
}
