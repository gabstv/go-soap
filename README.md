go-soap
=======

Send and receive data packed in SOAP 1.1  

Example:
```go
package sample

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
	v, err := soap.Marshal(msg)
	if err != nil {
		panic(err)
	}
	resp, err := v.Post("http://www.example.com/api/v1/soapwebservice.asmx")
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
