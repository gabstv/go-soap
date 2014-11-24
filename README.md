go-soap
=======

### Send and receive data wrapped in SOAP v1.1  
HTTP, SMTP, TCP, UDP or JMS  

Example:
```go
package main

import(
	"github.com/gabstv/go-soap"
	"fmt"
)

type MyMessage struct {
	Foo string
	Bar MBar
}

type MBar struct {
	Val int
	Baz string
}

type Resp struct {
	Success bool
	Error string
}

func main(){
	msg := &MyMessage{}
	msg.Foo = "test"
	msg.Bar.Val = 100
	msg.Bar.Baz = "foobarbaz"

	// soap.Marshal creates a *soap.Envelope that is ready to be submitted.
	v, err := soap.Marshal(msg)
	if err != nil {
		panic(err)
	}
	resp, err := v.Post("https://www.example.com/api/v1/soapwebservice.asmx")
	if err != nil {
		panic(err)
	}

	out := &Resp{}
	err = resp.Unmarshal(resp)

	if err != nil {
		panic(err)
	}

	fmt.Println("Response:", out)
}
```  

###If you need to change the Content-Type:

```go
// after running soap.Marshal(content)
v.HTMLContentType = "application/xml+soap"
```

If you're not using HTML 1.1 to submit, you can also write to a buffer:  

```go
err := v.WriteTo(some_io_Writer)
```
