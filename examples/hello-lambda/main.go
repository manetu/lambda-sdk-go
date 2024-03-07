package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gitlab.com/manetu/pre-release/lambda/lambda-sdk-go"
	"os"
)

type context struct {
}

func (c context) Handler(request lambda.Request) lambda.Response {

	log.Printf("handling request %v", request.Params)

	greeting := fmt.Sprintf("Hello, %s", request.Params["name"])

	return lambda.Response{
		Status:  200,
		Headers: lambda.Headers{"Content-Type": "text/plain"},
		Body:    greeting}
}

func main() {
	lambda.Init(&context{})
	log.Print("Module initialized:", os.Environ())
}
