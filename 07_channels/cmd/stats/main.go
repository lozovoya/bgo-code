package main

import (
	"fmt"
	"lectionchannels/pkg/stats"
	"log"
	"os"
	"runtime/trace"
	"sync"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()
	err = trace.Start(f)
	if err != nil {
		log.Fatal(err)
	}
	defer trace.Stop()

	const users = 10_000_001
	const transactionsPerUser = 100
	const transactionAmount = 1_00
	transactions := make([]int64, users * transactionsPerUser)
	for index := range transactions {
		// для простоты храним только суммы
		// и считаем, что каждая транзакция = 1 руб.
		transactions[index] = transactionAmount
	}

	total := int64(0)
	const partsCount = 10
	partSize := len(transactions) / partsCount
	wg := sync.WaitGroup{}
	wg.Add(partsCount)
	mu := sync.Mutex{}
	for i := 0; i < partsCount; i++ {
		part := transactions[i * partSize : (i + 1) * partSize]
		go func() {
			sum := stats.Sum(part)
			mu.Lock()
			total += sum
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()

	mu.Lock()
	fmt.Println(total)
	mu.Unlock()
}
