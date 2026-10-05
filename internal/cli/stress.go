package cli

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

type config struct {
	URL         string
	Requests    int
	Concurrency int
}

type report struct {
	URL             string
	Elapsed         time.Duration
	TotalRequests   int
	StatusCounts    map[int]int
	TransportErrors int
}

type result struct {
	statusCode int
	err        error
}

func runLoadTest(cfg config) report {
	workerCount := cfg.Concurrency
	if workerCount > cfg.Requests {
		workerCount = cfg.Requests
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	jobs := make(chan struct{})
	results := make(chan result)
	startedAt := time.Now()

	for range workerCount {
		go func() {
			for range jobs {
				response, err := client.Get(cfg.URL)
				if err != nil {
					results <- result{err: err}
					continue
				}
				_, _ = io.Copy(io.Discard, response.Body)
				_ = response.Body.Close()
				results <- result{statusCode: response.StatusCode}
			}
		}()
	}

	go func() {
		for range cfg.Requests {
			jobs <- struct{}{}
		}
		close(jobs)
	}()

	statusCounts := make(map[int]int)
	transportErrors := 0
	for range cfg.Requests {
		requestResult := <-results
		if requestResult.err != nil {
			transportErrors++
			continue
		}
		statusCounts[requestResult.statusCode]++
	}

	return report{
		URL:             cfg.URL,
		Elapsed:         time.Since(startedAt),
		TotalRequests:   cfg.Requests,
		StatusCounts:    statusCounts,
		TransportErrors: transportErrors,
	}
}

func printReport(output io.Writer, executionReport report) {
	fmt.Fprintln(output, "=== Relatório do teste de carga ===")
	fmt.Fprintf(output, "URL: %s\n", executionReport.URL)
	fmt.Fprintf(output, "Tempo total: %s\n", executionReport.Elapsed.Round(time.Millisecond))
	fmt.Fprintf(output, "Total de requests realizados: %d\n", executionReport.TotalRequests)
	fmt.Fprintf(output, "Requests com status HTTP 200: %d\n", executionReport.StatusCounts[http.StatusOK])
	fmt.Fprintln(output, "Distribuição dos outros status HTTP:")

	statusCodes := make([]int, 0, len(executionReport.StatusCounts))
	for statusCode := range executionReport.StatusCounts {
		if statusCode != http.StatusOK {
			statusCodes = append(statusCodes, statusCode)
		}
	}
	sort.Ints(statusCodes)
	if len(statusCodes) == 0 {
		fmt.Fprintln(output, "  Nenhum")
	}
	for _, statusCode := range statusCodes {
		fmt.Fprintf(output, "  HTTP %d: %d\n", statusCode, executionReport.StatusCounts[statusCode])
	}
	fmt.Fprintf(output, "Falhas de conexão/timeout: %d\n", executionReport.TransportErrors)
}
