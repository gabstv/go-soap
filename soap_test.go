package soap

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type Tr struct {
	Foo string
	Bar int
	Baz bool
}

type tht struct {
	t *testing.T
}

func (t *tht) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// decode
	env := &Envelope{}

	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		t.t.Fatal(err)
	}

	err = xml.Unmarshal(b, env)
	if err != nil {
		t.t.Fatal(err)
	}
	dat := &Tr{}

	err = env.Unmarshal(dat)

	if err != nil {
		t.t.Fatal(err)
	}

	if dat.Foo != "foo" {
		t.t.Fatal("dat.Foo != foo")
	}

	dat.Foo = "foo"
	dat.Bar = 2
	dat.Baz = true

	env, err = Marshal(dat)

	if err != nil {
		t.t.Fatal(err)
	}

	err = env.WriteTo(w)

	if err != nil {
		t.t.Fatal(err)
	}
}

func TestSoap(t *testing.T) {
	//
	h := &tht{t}
	go func() {
		err := http.ListenAndServe("127.0.0.1:13923", h)
		if err != nil {
			t.Fatal("Could not serve.", err)
		}
	}()
	time.Sleep(time.Millisecond * 100)

	dat := &Tr{"foo", 1, false}

	env, err := Marshal(dat)

	if err != nil {
		t.Fatal(err)
	}

	env, err = env.Post("http://localhost:13923/")

	if err != nil {
		t.Fatal(err)
	}

	err = env.Unmarshal(dat)

	if err != nil {
		t.Fatal(err)
	}

	if dat.Bar != 2 {
		t.Fatal("Bar != 2")
	}

	if !dat.Baz {
		t.Fatal("Baz != true")
	}

	//
}

func TestUnmarshal(t *testing.T) {
	vvv := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AuthorizeTransaction xmlns="https://www.example.com/webservice/pay">
      <request>
        <OrderData>
          <MerchantId>MI5678</MerchantId>
          <OrderId>OID3949v1</OrderId>
        </OrderData>
      </request>
    </AuthorizeTransaction>
  </soap:Body>
</soap:Envelope>`
	vv := struct {
		XMLName xml.Name `xml:"AuthorizeTransaction"`
		Request struct {
			OrderData struct {
				MerchantId string
				OrderId    string
			}
		} `xml:"request"`
	}{}
	ev := Envelope{}
	err := xml.Unmarshal([]byte(vvv), &ev)

	if err != nil {
		t.Fatal(err)
	}

	err = ev.Unmarshal(&vv)

	if err != nil {
		t.Fatal(err)
	}

	if vv.Request.OrderData.MerchantId != "MI5678" {
		t.Fatal("Unmarshal error:", vv)
	}

	if vv.Request.OrderData.OrderId != "OID3949v1" {
		t.Fatal("Unmarshal error:", vv)
	}
}
