package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mustafaakin/lambda-performance-test"
	"github.com/getlantern/errors"
)

type Request struct {
	Iterations int
	Batch      int
	Threads    int
}

func Handler(request Request) (*lambda_performance_test.TestResults, error) {
	warmup := lambda_performance_test.DoLongWork(10, 5, 2)
	if warmup.Mean == 0 {
		// To ensure it is not optimized out
		return nil, errors.New("no mean")
	}

	result := lambda_performance_test.DoLongWork(request.Iterations, request.Batch, request.Threads)
	return &result, nil
}

func main() {
	lambda.Start(Handler)
}
