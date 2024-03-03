package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"sync"

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

	if err := executeAndProfile(file, &wordStruct); err != nil {
		fmt.Println("Error processing file:", err)
	}
}

func executeAndProfile(file *os.File, wordStruct *WordStruct) error {
	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup

	for scanner.Scan() {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()

			wordStruct.mut.Lock()
			defer wordStruct.mut.Unlock()

			l := strings.Split(line, ";")

			if len(l) > 1 {
				word := l[0]
				count := l[1]

				decimalCount, err := decimal.NewFromString(count)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				item := wordStruct.hashmap[word]

				//max
				if decimalCount.GreaterThan(wordStruct.hashmap[word].Max) {
					item.Max = decimalCount
					wordStruct.hashmap[word] = item
				}

				//min
				if decimalCount.LessThan(wordStruct.hashmap[word].Min) {
					item.Min = decimalCount
					wordStruct.hashmap[word] = item
				}

				roundedDecimal := decimalCount.Round(2)

				item.WordCount = wordStruct.hashmap[word].WordCount + 1
				item.Amount = wordStruct.hashmap[word].Amount.Add(roundedDecimal)
				wordStruct.hashmap[word] = item
			}

		}(scanner.Text())
	}

	wg.Wait()

	for k, v := range wordStruct.hashmap {
		div := v.Amount.Div(decimal.NewFromInt(int64(v.WordCount))).Round(2)
		fmt.Println(k, div.String(), "count:", v.WordCount)
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return nil
}
