package sparql

import (
	"encoding/json"
	"errors"
	"fmt"
	core "github.com/manetu/lambda-sdk-go"
	"github.com/manetu/lambda-sdk-go/internal"
)

type Binding struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Solution map[string]Binding

type InnerResults struct {
	Bindings []Solution `json:"bindings"`
}

type Results struct {
	Results InnerResults `json:"results"`
}

type Response struct {
	Status  int          `json:"status"`
	Headers core.Headers `json:"headers"`
	Body    string       `json:"body"`
}

func Query(expr string) (*Results, error) {
	resp := internal.ManetuLambda0_0_1_SparqlQuery(expr)

	var response Response
	err := json.Unmarshal([]byte(resp), &response)
	if err != nil {
		fmt.Printf("unmarshal-error: %v", err)
		return nil, err
	}

	if response.Status != 200 {
		return nil, errors.New("query failure")
	}

	if response.Headers["content-type"] != "application/json" {
		return nil, errors.New("unexpected content type")
	}

	var result Results
	err = json.Unmarshal([]byte(response.Body), &result)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	return &result, nil
}
