# aws-lambda-middleware-example

**From**: a simple app that knows to say "i'm an app" in terminal

**To**: an AWS Lambda written in Go that has a logger middleware.

<details>
  <summary><b>1. Print to terminal</b></summary>
  We have a small app that print "i'm a app" to terminal.

```go
package main

import "fmt"

func handler() {
	fmt.Println("i'm a cli app")
}

func main() {
	fmt.Println("app started")
	handler()
}

//=================================
// > go run main.go
// started
// i'm a app
//=================================
```
</details>

<details>
  <summary><b>2. Transform in a web app</b></summary>
  The app is running and does the job, now we want to transform it into a web app.
In web apps we usually log the requests, so a new func is introduced.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// the extra set of instructions
// things to be done before running the business logic
func logging(r *http.Request) {
	log.Printf("remote_addr: %s", r.RemoteAddr)
}

// the business logic
// we could imagine here is defined the codebase core capability
func handler(w http.ResponseWriter, r *http.Request) {
	logging(r)
	fmt.Fprintf(w, "i'm a web app")
}

func main() {
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// start the web server
//=================================
// > go run main.go
// 2021/02/24 15:16:14 remote_addr: 127.0.0.1:47834
//=================================

// call it 
//=================================
// > curl http://localhost:8080
// i'm an web app%
//=================================
```
</details>

<details>
  <summary><b>3. Introduce middleware for decoupling</b></summary>
  We know that the web ap does its job, but we want more than that: decoupled code.
This will translate in transforming the `logging` func in a middleware.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// the extra set of instructions
// things to be done before running the business logic
func logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("remote_addr: %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

// the business logic
// we could imagine here is defined the codebase core capability
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "i'm a web app")
}

func main() {
	http.HandleFunc("/", logging(handler))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// start the web server
//=================================
// > go run main.go
// 2021/02/24 15:16:15 remote_addr: 127.0.0.1:47835
//=================================
// call it 
//=================================
// > curl http://localhost:8080
// i'm an web app%
//=================================
```

Same output, transforming `logging` func in a middleware works like a charm.

 
</details>

<details>
  <summary><b>4. Running on AWS Lambda - without middleware</b></summary>
  Now is time for production. It is a small and straightforward web app, so we take the decision to hosted via AWS Lambda.

First we will setup the stage for running on AWS Lambda the business logic.

```go
package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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
	lambda.Start(handle)
}
```
For deploying and testing follow the AWS official [docs](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html).
</details>

<details>
  <summary><b>5. Running on AWS Lambda - put back the middleware</b></summary>
  Put back the `logging` func adapted for AWS Lambda:

  ```go
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
```

</details>
