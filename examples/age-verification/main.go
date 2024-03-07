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
		log.Err(err)
		return lambda.Response{
			Status: 500,
		}
	}

	switch len(r.Results.Bindings) {
	case 0:
		return lambda.Response{
			Status:  200,
			Headers: lambda.Headers{"Content-Type": "text/plain"},
			Body:    "not-found",
		}
	case 1:
		break
	default:
		return lambda.Response{
			Status: 500,
			Body:   "unexpected multiple matching results",
		}
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
