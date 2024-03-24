package lambda

// #include <stdlib.h>
import "C"
import (
	"encoding/json"
	"github.com/manetu/lambda-sdk-go/internal"
	"unsafe"
)

type Headers map[string]string
type Params map[string]any

type Request struct {
	Headers  Headers `json:"headers"`
	PathInfo string  `json:"path-info"`
	Params   Params  `json:"params"`
}

type Response struct {
	Status  int     `json:"status"`
	Headers Headers `json:"headers"`
	Body    any     `json:"body"`
}

type Lambda interface {
	Handler(Request) Response
}

type context struct {
	lambda Lambda
}

func (c context) HandleRequest(request string) string {
	var req Request
	err := json.Unmarshal([]byte(request), &req)
	if err != nil {
		v, err := json.Marshal(&Response{Status: 500})
		if err != nil {
			panic(err)
		}
		return string(v)
	}

	resp := c.lambda.Handler(req)

	v, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return string(v)
}

func (c context) Malloc(len uint32) uint32 {
	return uint32(uintptr(C.malloc(C.ulong(len))))
}

func (c context) Free(ptr uint32) {
	C.free(unsafe.Pointer(uintptr(ptr)))
}

func Init(lambda Lambda) {
	internal.SetExportsManetuLambda0_0_2_Guest(&context{lambda: lambda})
}
