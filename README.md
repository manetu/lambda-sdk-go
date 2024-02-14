# lambda-sdk-go

This repository hosts an SDK for developing Lambda functions for the Manetu Platform in the [go](https://go.dev/) programming language.

## Prerequisites

- [tinygo](https://tinygo.org/)
- [wasm-to-oci](https://github.com/engineerd/wasm-to-oci)

## Project setup

### Create a new go project

``` shell
$ go mod init my-lambda
```

### Install dependencies

This SDK may be installed with

``` shell
$ go get -u gitlab.com/manetu/pre-release/lambda/lambda-sdk-go
```

We will also be using [zerolog](https://github.com/rs/zerolog) and [mustach](https://github.com/hoisie/mustache) as part of our example

``` shell
$ go get -u github.com/rs/zerolog/log
$ go get -u github.com/hoisie/mustache
```

### Create a main module and HTTP event handler

Create a file 'main.go' with the following contents:

``` golang
package main

import (
	"github.com/hoisie/mustache"
	"github.com/rs/zerolog/log"
	"gitlab.com/manetu/pre-release/lambda/lambda-sdk-go"
	"gitlab.com/manetu/pre-release/lambda/lambda-sdk-go/sparql"
	"os"
	"time"
)

type context struct {
}

var queryTemplate = `
PREFIX foaf:  <http://xmlns.com/foaf/0.1/> 
SELECT ?dob
WHERE { 
     ?s foaf:biometric-hash "{{biometric-hash}}" ;
        foaf:dob            ?dob .
}
`

func (c context) Handler(request lambda.Request) lambda.Response {

	log.Printf("handling request %v", request.Params)

	query := mustache.Render(queryTemplate, request.Params)
	r, err := sparql.Query(query)
	if err != nil {
		panic(err)
	}

	if len(r.Results.Bindings) != 1 {
		panic("too many results")
	}

	dob := r.Results.Bindings[0]["dob"]

	if dob.Type != "xsd:date" {
		panic("unexpected type")
	}

	now := time.Now()
	// 8760 hours in a 365 day year, * 21 years = 183960
	minimum := now.Add(time.Duration(-183960) * time.Hour)
	d, err := time.Parse(time.DateOnly, dob.Value)

	return lambda.Response{
		Status:  200,
		Headers: lambda.Headers{"Content-Type": "text/plain"},
		Body:    d.Before(minimum)}
}

func main() {
	lambda.Init(&context{})
	log.Print("Module initialized:", os.Environ())
}
```

### Compile the program

The Manetu platform serves Lambda functions within a [WebAssembly](https://webassembly.org/) environment.  We can leverage the [WASI support](https://tinygo.org/docs/guides/webassembly/wasi/) in tinygo to compile our program.

``` shell
$ tinygo build -o my-lambda.wasm --target=wasi main.go
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
  name: verify
  tag: verification
spec:
  runtime: wasi.1.alpha1
  image: oci://my-registry.example.com/my-lambda:v0.0.1
  env:
    LOG_LEVEL: trace
  permissions:
    assumed-roles:
      - mrn:iam:manetu.io:role:admin
    scopes:
      - mrn:iam:manetu.io:scope:read-api
  triggers:
    http-queries:
      - route: /over21
        summary: "Verify a person is at least 21 years old by biometric hash"
        description: "This request allows you to determine if the person identified by a biometric hash is at least 21 years old as of the time the call is made."
        query-parameters:
          - name: "biometric-hash"
            schema: { type: "string" }
            description: "The biometric hash of the user"
        responses:
          200:
            description: "verification result when the biometric hash is valid"
            content:
              text/plain:
                schema:
                  type: boolean
          404:
            description: "the biometric hash is not found"
```

Be sure to adjust the image OCI url
