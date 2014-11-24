package soap

import "encoding/xml"
import "net/http"
import "bytes"
import "io"

type M map[string]string

type Envelope struct {
	XMLName         xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Xsi             string   `xml:"xmlns:xsi,attr"`
	Soapenc         string   `xml:"xmlns:soapenc,attr"`
	Xsd             string   `xml:"xmlns:xsd,attr"`
	EncodingStyle   string   `xml:"soap:encodingStyle,attr"`
	Soap            string   `xml:"xmlns:soap,attr"`
	Body            Body
	HTTPContentType string `xml:"-"`
}

type Body struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Data    string   `xml:",innerxml"`
}

func Marshal(data interface{}) (*Envelope, error) {
	msg, err := xml.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		Xsi:             "http://www.w3.org/2001/XMLSchema-instance",
		Soapenc:         "http://schemas.xmlsoap.org/soap/encoding/",
		Xsd:             "http://www.w3.org/2001/XMLSchema",
		EncodingStyle:   "http://schemas.xmlsoap.org/soap/encoding/",
		Soap:            "http://schemas.xmlsoap.org/soap/envelope/",
		Body:            Body{Data: string(msg)},
		HTTPContentType: "text/xml; charset=utf-8", // it could also be application/soap+xml
	}, nil
}

func (env *Envelope) WriteTo(writer io.Writer) error {
	msg, err := xml.Marshal(env)
	if err != nil {
		return err
	}

	writer.Write(msg)
	return nil
}

func (env *Envelope) Post(url string) (response *Envelope, err error) {
	buf := new(bytes.Buffer)
	buf.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	err = env.WriteTo(buf)

	if err != nil {
		return
	}

	var htresp *http.Response
	htresp, err = http.Post(url, env.HTTPContentType, buf)

	if err != nil {
		return
	}

	defer htresp.Body.Close()

	dec := xml.NewDecoder(htresp.Body)
	response = &Envelope{}
	err = dec.Decode(response)

	return
}

func (env *Envelope) PostAdv(url string, headers M) (response *Envelope, err error) {
	buf := new(bytes.Buffer)
	buf.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	err = env.WriteTo(buf)

	if err != nil {
		return
	}

	var htresp *http.Response

	req, err := http.NewRequest("POST", url, buf)

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", env.HTTPContentType)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	htresp, err = http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer htresp.Body.Close()

	dec := xml.NewDecoder(htresp.Body)
	response = &Envelope{}
	err = dec.Decode(response)

	return
}

func (env *Envelope) Unmarshal(target interface{}) error {
	return xml.Unmarshal([]byte(env.Body.Data), target)
}
