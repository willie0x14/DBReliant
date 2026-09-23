package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func main() {
	requests := flag.Int("requests", 100, "total number of transfer requests")
	workers := flag.Int("workers", 20, "number of concurrent workers")
	amount := flag.Int64("amount", 1, "transfer amount")
	endpoint := flag.String(
		"endpoint",
		"http://localhost:8080/transfers",
		"transfer API endpoint",
	)

	flag.Parse()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	jobs := make(chan int, *requests)

	var success atomic.Int64
	var failed atomic.Int64

	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < *workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for job := range jobs {
				req := buildRequest(job, *amount)

				if err := sendTransfer(client, *endpoint, req); err != nil {
					failed.Add(1)
					continue
				}

				success.Add(1)
			}
		}()
	}

	for i := 0; i < *requests; i++ {
		jobs <- i
	}

	close(jobs)

	wg.Wait()

	duration := time.Since(start)

	fmt.Printf("requests: %d\n", *requests)
	fmt.Printf("workers: %d\n", *workers)
	fmt.Printf("success: %d\n", success.Load())
	fmt.Printf("failed: %d\n", failed.Load())
	fmt.Printf("duration: %s\n", duration)
}

func buildRequest(job int, amount int64) transferRequest {
	if job%2 == 0 {
		return transferRequest{
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        amount,
		}
	}

	return transferRequest{
		FromAccountID: 2,
		ToAccountID:   1,
		Amount:        amount,
	}
}

func sendTransfer(
	client *http.Client,
	endpoint string,
	req transferRequest,
) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return nil
}
