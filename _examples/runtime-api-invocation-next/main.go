package main

import (
	"context"
	"fmt"
	"log"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/michimani/aws-lambda-api-go/alago"
	alagoruntime "github.com/michimani/aws-lambda-api-go/runtime"
)

type Response struct {
	Message string `json:"message"`
}

type ExampleEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func handleRequest(ctx context.Context, event ExampleEvent) (*Response, error) {
	log.Println("start handler")
	defer log.Println("end handler")

	log.Printf("event: %#+v", event)

	// Get AWS Request ID from Lambda Runtime API.
	ac, err := alago.NewClient(&alago.NewClientInput{})
	if err != nil {
		return nil, err
	}

	out, err := alagoruntime.InvocationNext(context.Background(), ac)
	if err != nil {
		return nil, err
	}

	log.Printf("%#+v", out)

	return &Response{
		Message: fmt.Sprintf("Request ID is %s", out.AWSRequestID),
	}, nil
}

func init() {
	log.Println("cold start")
}

func main() {
	runtime.Start(handleRequest)
}
