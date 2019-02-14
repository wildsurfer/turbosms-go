package wsdl

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

var cookieJar, _ = cookiejar.New(nil)
var Debug bool

type Auth struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo Auth"`

	Login    string `xml:"login,omitempty"`
	Password string `xml:"password,omitempty"`
}

type AuthResponse struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo AuthResponse"`

	AuthResult string `xml:"AuthResult,omitempty"`
}

type GetCreditBalance struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo GetCreditBalance"`
}

type GetCreditBalanceResponse struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo GetCreditBalanceResponse"`

	GetCreditBalanceResult string `xml:"GetCreditBalanceResult,omitempty"`
}

type SendSMS struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo SendSMS"`

	Sender      string `xml:"sender,omitempty"`
	Destination string `xml:"destination,omitempty"`
	Text        string `xml:"text,omitempty"`
	Wappush     string `xml:"wappush,omitempty"`
}

type SendSMSResponse struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo SendSMSResponse"`

	SendSMSResult *ArrayOfString `xml:"SendSMSResult,omitempty"`
}

type GetNewMessages struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo GetNewMessages"`
}

type GetNewMessagesResponse struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo GetNewMessagesResponse"`

	GetNewMessagesResult *ArrayOfString `xml:"GetNewMessagesResult,omitempty"`
}

type GetMessageStatus struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo GetMessageStatus"`

	MessageId string `xml:"MessageId,omitempty"`
}

type GetMessageStatusResponse struct {
	XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo GetMessageStatusResponse"`

	GetMessageStatusResult string `xml:"GetMessageStatusResult,omitempty"`
}

type ArrayOfString struct {
	//XMLName xml.Name `xml:"http://turbosms.in.ua/api/Turbo ArrayOfString"`

	ResultArray []string `xml:"ResultArray,omitempty"`
}

type ServiceSoap struct {
	client *SOAPClient
}

func NewServiceSoap(url string, tls bool, auth *BasicAuth) *ServiceSoap {
	if url == "" {
		url = "http://turbosms.in.ua/api/soap.html"
	}
	client := NewSOAPClient(url, tls, auth)

	return &ServiceSoap{
		client: client,
	}
}

func (service *ServiceSoap) Auth(request *Auth) (*AuthResponse, error) {
	response := new(AuthResponse)
	err := service.client.Call("http://turbosms.in.ua/api/Turbo/Auth", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *ServiceSoap) GetCreditBalance(request *GetCreditBalance) (*GetCreditBalanceResponse, error) {
	response := new(GetCreditBalanceResponse)
	err := service.client.Call("http://turbosms.in.ua/api/Turbo/GetCreditBalance", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *ServiceSoap) SendSMS(request *SendSMS) (*SendSMSResponse, error) {
	response := new(SendSMSResponse)
	err := service.client.Call("http://turbosms.in.ua/api/Turbo/SendSMS", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *ServiceSoap) GetNewMessages(request *GetNewMessages) (*GetNewMessagesResponse, error) {
	response := new(GetNewMessagesResponse)
	err := service.client.Call("http://turbosms.in.ua/api/Turbo/GetNewMessages", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *ServiceSoap) GetMessageStatus(request *GetMessageStatus) (*GetMessageStatusResponse, error) {
	response := new(GetMessageStatusResponse)
	err := service.client.Call("http://turbosms.in.ua/api/Turbo/GetMessageStatus", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Body SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Header interface{}
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

type BasicAuth struct {
	Login    string
	Password string
}

type SOAPClient struct {
	url  string
	tls  bool
	auth *BasicAuth
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, tls bool, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url:  url,
		tls:  tls,
		auth: auth,
	}
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{
		//Header:        SoapHeader{},
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}
	if Debug == true {
		log.Println(buffer.String())
	}

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}
	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	if soapAction != "" {
		req.Header.Add("SOAPAction", soapAction)
	}

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: s.tls,
		},
		Dial: dialTimeout,
	}

	client := &http.Client{Transport: tr, Jar: cookieJar}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		log.Println("empty response")
		return nil
	}
	if Debug == true {
		log.Println(string(rawbody))
	}
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}
