package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := "postgres://dbreliant_app:labpass@localhost:5432/dbreliant?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	const workers = 100

	var wg sync.WaitGroup
	wg.Add(workers)

	start := time.Now()
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()

			// _, err := db.Exec("SELECT pg_sleep(2)")
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := db.ExecContext(ctx, "SELECT pg_sleep(2)")
			if err != nil {
				fmt.Printf("worker %d error: %v\n", id, err)
			}
		}(i)
	}
	time.Sleep(500 * time.Millisecond)

	stats := db.Stats()

	fmt.Println("=== DURING LOAD ===")
	fmt.Println("OpenConnections:", stats.OpenConnections)
	fmt.Println("InUse:", stats.InUse)
	fmt.Println("Idle:", stats.Idle)
	fmt.Println("WaitCount:", stats.WaitCount)
	fmt.Println("WaitDuration:", stats.WaitDuration)
	fmt.Println("MaxIdleClosed:", stats.MaxIdleClosed)

	wg.Wait()

	workloadDuration := time.Since(start)

	time.Sleep(1 * time.Second)

	finalStats := db.Stats()

	fmt.Println("=== AFTER LOAD ===")
	fmt.Println("OpenConnections:", finalStats.OpenConnections)
	fmt.Println("InUse:", finalStats.InUse)
	fmt.Println("Idle:", finalStats.Idle)
	fmt.Println("WaitCount:", finalStats.WaitCount)
	fmt.Println("WaitDuration:", finalStats.WaitDuration)
	fmt.Println("MaxIdleClosed:", finalStats.MaxIdleClosed)
	fmt.Println("Workload duration:", workloadDuration)

	fmt.Println("Total duration:", time.Since(start))

}
