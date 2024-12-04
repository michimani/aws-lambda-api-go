package main

import (
	"context"
	"log"

	runtime "github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func handleRequest(ctx context.Context) (*Response, error) {
	log.Println("start handler")
	defer log.Println("end handler")

	return &Response{
		Message: "Hello Lambda!",
	}, nil
}

func init() {
	log.Println("cold start")
}

func main() {
	runtime.Start(handleRequest)
}
