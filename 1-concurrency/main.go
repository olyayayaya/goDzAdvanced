package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {

	numsSlice := []int{}
	for i := 0; i < 10; i++ {
		numsSlice = append(numsSlice, rand.Intn(100))
	}
	
	numsCh := make(chan int, len(numsSlice))
	resultCh := make(chan int, len(numsSlice))
	
	for _, num := range numsSlice {
		numsCh <- num
	}
	
	close(numsCh)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go sq(numsCh, resultCh, &wg)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

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
