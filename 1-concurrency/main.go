package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {

	numsCh := make(chan int, 10)
	resultCh := make(chan int, 10)

	go fullSlice(numsCh)

	var wg sync.WaitGroup

	wg.Add(1)
	go sq(numsCh, resultCh, &wg)

	wg.Wait()
	close(resultCh)

	resultSlice := []int{}
	for res := range resultCh {
		resultSlice = append(resultSlice, res)
	}

	fmt.Println(resultSlice)
}

func sq(numsCh chan int, resCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range numsCh {
		resCh <- num * num
	}
}

func fullSlice(numsCh chan int) {
	numsSlice := []int{}
	for i := 0; i < 10; i++ {
		numsSlice = append(numsSlice, rand.Intn(100))
	}

	for _, num := range numsSlice {
		numsCh <- num
	}
	close(numsCh)
}
