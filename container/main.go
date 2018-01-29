package main

import (
	"github.com/mustafaakin/lambda-performance-test"
	"encoding/json"
	"os"
	"flag"
	"log"
)

var (
	iterations = flag.Int("iterations", 10, "The number of iterations")
	batch      = flag.Int("batch", 50, "The batch size")
	threads    = flag.Int("threads", 1, "The thread size")
)

func main() {
	flag.Parse()

	// The warmup
	warmup := lambda_performance_test.DoLongWork(10, 5, 2)
	if warmup.Mean == 0 {
		// To ensure it is not optimized out
		log.Fatal("No mean")
	}

	result := lambda_performance_test.DoLongWork(*iterations, *batch, *threads)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(result)
}
