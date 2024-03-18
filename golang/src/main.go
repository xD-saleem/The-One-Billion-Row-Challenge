package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/alphadose/haxmap"
)

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

	file, err := os.Open("../../1bill.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	result := processFile(file)
	fmt.Println(result)
}

type result struct {
	city string
	min  float64
	max  float64
	avg  float64
}

type info struct {
	count int64
	min   float64
	max   float64
	sum   float64
}

func processFile(input *os.File) string {
	m, err := readChunkByChunk(input)
	if err != nil {
		panic(err)
	}

	resultArr := make([]result, m.Len())
	var count int

	m.ForEach(func(city string, value *info) bool {
		resultArr[count] = result{
			city: city,
			min:  round(float64(value.min) / 10.0),
			max:  round(float64(value.max) / 10.0),
			avg:  round(float64(value.sum) / 10.0 / float64(value.count)),
		}
		count++
		return true
	})

	sort.Slice(resultArr, func(i, j int) bool {
		return resultArr[i].city < resultArr[j].city
	})

	var stringsBuilder strings.Builder
	for _, i := range resultArr {
		stringsBuilder.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", i.city, i.min, i.avg, i.max))
	}
	return stringsBuilder.String()[:stringsBuilder.Len()-2]
}

func readChunkByChunk(file *os.File) (*haxmap.Map[string, *info], error) {
	const chunkSize = 64 * 1024 * 1024 // 64MB

	m := haxmap.New[string, *info]()

	resultStream := make(chan map[string]info, 10)
	chunkStream := make(chan []byte, 15)

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU()-1; i++ {
		wg.Add(1)
		go func() {
			for chunk := range chunkStream {
				processChunk(chunk, resultStream)
			}
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, chunkSize)
		leftover := make([]byte, 0, chunkSize)
		for {
			readTotal, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			buf = buf[:readTotal]

			lastNewLineIndex := bytes.LastIndex(buf, []byte{'\n'})

			buf = append(leftover, buf[:lastNewLineIndex+1]...)
			leftover = make([]byte, len(buf[lastNewLineIndex+1:]))
			copy(leftover, buf[lastNewLineIndex+1:])

			chunkStream <- buf
		}

		close(chunkStream)
	}()

	go func() {
		wg.Wait()
		close(resultStream)
	}()

	for t := range resultStream {
		for city, info := range t {
			if val, ok := m.Get(city); ok {
				val.count += info.count
				val.sum += info.sum
				if info.min < val.min {
					val.min = info.min
				}

				if info.max > val.max {
					val.max = info.max
				}
				m.Set(city, val)
			} else {
				m.Set(city, &info)
			}
		}
	}

	return m, nil
}

func processChunk(buf []byte, resultStream chan<- map[string]info) {
	toSend := make(map[string]info)
	var start int
	var city string

	stringBuf := string(buf)

	for index, char := range stringBuf {
		switch char {
		case ';':
			city = stringBuf[start:index]
			start = index + 1
		case '\n':
			if (index-start) > 1 && len(city) != 0 {
				temp, _ := strconv.ParseFloat(stringBuf[start:index], 64)

				start = index + 1

				if val, ok := toSend[city]; ok {
					val.count++
					val.sum += temp
					if temp < val.min {
						val.min = temp
					}

					if temp > val.max {
						val.max = temp
					}
					toSend[city] = val
				} else {
					toSend[city] = info{
						count: 1,
						min:   temp,
						max:   temp,
						sum:   temp,
					}
				}

				city = ""
			}
		}
	}
	resultStream <- toSend
}

func round(x float64) float64 {
	rounded := math.Round(x * 10)
	if rounded == -0.0 {
		return 0.0
	}
	return rounded / 10
}
