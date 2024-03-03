package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
	"sync/atomic"

	"github.com/shopspring/decimal"
)

type WordCount struct {
	Amount    decimal.Decimal
	WordCount int
	Max       decimal.Decimal
	Min       decimal.Decimal
}

type WordStruct struct {
	hashmap map[string]WordCount
	mut     sync.Mutex
}

func main() {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	f, err = os.Create("mem.pprof")
	if err != nil {
		panic(err)
	}
	pprof.WriteHeapProfile(f)
	defer f.Close()

	wordStruct := WordStruct{
		hashmap: make(map[string]WordCount),
		mut:     sync.Mutex{},
	}

	file, err := os.Open("../../1bill.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if err := profileFileReadingBase(file, &wordStruct); err != nil {
		fmt.Println("Error processing file:", err)
	}
}

func profileFileReadingBase(file *os.File, wordStruct *WordStruct) error {
	var ch = make(chan string, 100)

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		for scanner.Scan() {
			line := scanner.Text()
			ch <- line
		}

		close(ch) // Close the channel inside the goroutine
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}
	}()

	go func() {
		defer wg.Done()
		processLines(ch, wordStruct) // Start a new goroutine for processing lines
	}()

	wg.Wait() // Wait for the scanner goroutine to finish
	println("done reading file")

	return nil
}

var counter int32

func processLines(ch <-chan string, wordStruct *WordStruct) {

	for range ch {
		atomic.AddInt32(&counter, 1)
	}
	fmt.Println("done processing lines ", counter)
}
