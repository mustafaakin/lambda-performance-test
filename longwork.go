package lambda_performance_test

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
	"log"
	"math"
	"sync"
	"runtime"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randomBytes(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func DoLongWork(totalSize, batchSize, threads int) TestResults {
	// We pre allocate for reducing the allocation cost at the iterations
	results := make([]int64, totalSize)

	for i := 0; i < totalSize; i++ {
		st := time.Now()

		var wg sync.WaitGroup
		for T := 0; T < threads; T++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for j := 1; j < batchSize; j++ {
					pass := randomBytes(64)
					bytes, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
					if err != nil || bytes == nil {
						log.Fatal("Error occurred: ", err)
					}
				}
			}()
		}

		wg.Wait()

		en := time.Now()
		results[i] = en.Sub(st).Nanoseconds() / 1e6
	}

	return prepareResults(results)
}

type TestResults struct {
	Durations              []int64
	Mean                   float64
	Variance               float64
	StdDev                 float64
	CoefficientOfVariation float64
	Max                    int64
	Min                    int64
	NoOfCores              int
}

func prepareResults(durations []int64) TestResults {
	tr := TestResults{Durations: durations}

	var total float64
	tr.Min = math.MaxInt64
	tr.Max = math.MinInt64
	for _, d := range durations {
		ms := d
		if ms > tr.Max {
			tr.Max = ms
		}
		if ms < tr.Min {
			tr.Min = ms
		}
		total += float64(ms)
	}

	tr.Mean = total / float64(len(durations))

	// Calculate variance
	tr.Variance = 0
	for _, d := range durations {
		tr.Variance += math.Pow(math.Abs(float64(d)-tr.Mean), 2)
	}

	tr.Variance = tr.Variance / float64(len(durations))
	tr.StdDev = math.Sqrt(tr.Variance)
	tr.CoefficientOfVariation = tr.StdDev / tr.Mean

	tr.NoOfCores = runtime.NumCPU()

	return tr
}
