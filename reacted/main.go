package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request context.Context) (events.LambdaFunctionURLResponse, error) {
	// events.LambdaFunctionURLRequest{}
	fmt.Println(request)

	greeting := "Hello World"

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       "{\"message\": \"" + greeting + "\"}",
	}, nil
}

func main() {
	lambda.Start(handler)
}
