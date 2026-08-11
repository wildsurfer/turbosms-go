package turbosms

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wildsurfer/turbosms-go/wsdl"
)

const sendSMSResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Body>
<SendSMSResponse xmlns="http://turbosms.in.ua/api/Turbo">
<SendSMSResult>
<ResultArray>Сообщения успешно отправлены</ResultArray>
<ResultArray>a3b0cfb6-0000-0000-0000-000000000000</ResultArray>
</SendSMSResult>
</SendSMSResponse>
</soap:Body>
</soap:Envelope>`

func TestSendSMS(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("SOAPAction") == "" {
			t.Error("missing SOAPAction header")
		}
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		w.Write([]byte(sendSMSResponse))
	}))
	defer ts.Close()

	c := &Client{SoapService: wsdl.NewServiceSoap(ts.URL, false, &wsdl.BasicAuth{})}
	r, err := c.SendSMS("sender", "+380000000000", "test", "")
	if err != nil {
		t.Fatal(err)
	}
	if r.SendSMSResult == nil || len(r.SendSMSResult.ResultArray) != 2 {
		t.Fatalf("unexpected result: %+v", r)
	}
}
