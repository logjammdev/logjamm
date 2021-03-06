package logjamm

import (
	"errors"
	"time"

	"github.com/logjammdev/utils"
	"github.com/mxschmitt/playwright-go"
)

type Test func() map[string]int64

func WebTest(options Options, test Test) []map[string]int64 {
	var results []map[string]int64
	if options.Iterations != 0 && options.Duration == 0 {
		results = webTestIterations(options.Iterations, test)
	} else if options.Iterations == 0 {
		results = webTestDuration(options.Duration, test)
	} else {
		err := errors.New("error options: duration and iteration cannot be used at the same time")
		utils.Explain(err)
	}
	return results
}

func webTestIterations(iterations int, f func() map[string]int64) []map[string]int64 {
	var results []map[string]int64
	for i := 0; i < iterations; i++ {
		results = append(results, f())
	}
	return results
}

func webTestDuration(duration int, f func() map[string]int64) []map[string]int64 {
	var results []map[string]int64
	after := time.Now().Add(time.Duration(duration) * time.Second)
	for {
		now := time.Now()
		results = append(results, f())
		if now.After(after) {
			break
		}
	}
	return results
}

func Step(label string, f func(page playwright.Page) playwright.Page, page playwright.Page) (string, int64, playwright.Page) {
	start := time.Now()
	page = f(page)
	return label, time.Since(start).Milliseconds(), page
}
