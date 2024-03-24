# lambda-sdk-go

This repository hosts an SDK for developing Lambda functions for the Manetu Platform in the [go](https://go.dev/) programming language.

## Prerequisites

- [tinygo](https://tinygo.org/)
- [wasm-to-oci](https://github.com/engineerd/wasm-to-oci)

## Project setup

### Create a directory for your project

``` shell
mkdir my-lambda
cd my-lambda
```

### Initialize a new go project in your working directory

``` shell
go mod init my-lambda
```

### Install dependencies

This SDK may be installed with

``` shell
go get -u github.com/manetu/lambda-sdk-go
```

We will also be using [zerolog](https://github.com/rs/zerolog) as part of our example

``` shell
go get -u github.com/rs/zerolog/log
```

### Create a main module and HTTP event handler

Create a file 'main.go' with the following contents:

``` golang
package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/manetu/lambda-sdk-go"
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

```

### Compile the program

The Manetu platform serves Lambda functions within a [WebAssembly](https://webassembly.org/) environment.  We can leverage the [WASI support](https://tinygo.org/docs/guides/webassembly/wasi/) in tinygo to compile our program.

``` shell
tinygo build -o my-lambda.wasm --target=wasi main.go
```

### Publish the WASM code

We can leverage any [OCI](https://opencontainers.org/) registry to publish our Lambda function using the [wasm-to-oci](https://github.com/engineerd/wasm-to-oci) tool.

``` shell
$ wasm-to-oci push my-lambda.wasm my-registry.example.com/my-lambda:v0.0.1
INFO[0003] Pushed: my-registry.example.com/my-lambda:v0.0.1
INFO[0003] Size: 1242738
INFO[0003] Digest: sha256:cf9040f3bcd0e84232013ada2b3af02fe3799859480046b88cdd03b59987f3c9
```

### Define a specification for your Lambda function

Create a file 'site.yml' with the following contents:

``` yaml
api-version: lambda.manetu.io/v1alpha1
kind: Site
metadata:
  name: hello
spec:
  runtime: wasi.1.alpha2
  image: oci://my-registry.example.com/my-lambda:v0.0.1
  env:
    LOG_LEVEL: trace
  triggers:
    http-queries:
      - route: /greet
        summary: "Returns a greeting to the user"
        description: "This request allows you to test the ability to deploy and invoke a simple lambda function."
        query-parameters:
          - name: "name"
            schema: { type: "string" }
            description: "The caller's name"
        responses:
          200:
            description: "computed greeting"
            content:
              text/plain:
                schema:
                  type: string
```

Be sure to adjust the image OCI url
