package main

import (
	"fmt"
	"github.com/manetu/lambda-sdk-go"
	"github.com/rs/zerolog/log"
	"os"
)

type context struct {
}

func (c context) Handler(request lambda.Request) lambda.Response {
	return lambda.Response{
		Status:  200,
		Headers: lambda.Headers{"Content-Type": "text/plain"},
		Body:    fmt.Sprintf("Hello, %s", request.Params["name"])}
}

func main() {
	lambda.Init(&context{})
	log.Print("Module initialized:", os.Environ())
}
