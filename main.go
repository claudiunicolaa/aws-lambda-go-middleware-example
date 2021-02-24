package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

type handlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// the extra set of instructions
// things to be done before running the business logic
func logging(f handlerFunc) handlerFunc {
	return func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		response, err := f(ctx, r)
		log.Printf("remote_addr: %s", r.RequestContext.Identity.SourceIP)
		return response, err
	}
}

// the business logic
// we could imagine here is defined the codebase core capability
func handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "i'm an app",
		StatusCode: 200,
		Headers: map[string] string {
			"Content-Type":   "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(logging(handle) )
}