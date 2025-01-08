package main

import (
	"fmt"
	"sync"
)

func main() {
	dataChannel := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	go producer(dataChannel, &wg)
	go consumer(dataChannel, &wg)
	wg.Wait()
}

func producer(dataChannel chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(dataChannel)
	for i := 0; i < 10; i++ {
		dataChannel <- i
	}
}

func consumer(dataChannel <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for x := range dataChannel {
		fmt.Println(x)
	}
}
