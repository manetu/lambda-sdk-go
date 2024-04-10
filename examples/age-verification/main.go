package main

import (
	"github.com/hoisie/mustache"
	"github.com/manetu/lambda-sdk-go"
	"github.com/manetu/lambda-sdk-go/sparql"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type context struct {
}

var queryTemplate = `
PREFIX id: <http://example.gov/rmv/>
SELECT ?dob
WHERE {
     ?s id:biometric-hash "{{biometric-hash}}" ;
        id:dob            ?dob .
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
