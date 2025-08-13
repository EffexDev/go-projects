package main

import (
	"fmt"
	"sync"
)

func adder() func(a int) int {
	total := 0
	return func(x int) int {
		total += x
		return total
	}
}

func main() {

	sum := adder()

	//Start a worker pool
	var textWaitGroup sync.WaitGroup

	//Loop condition so the async eventually stops
	for i := 1; i <= 10; i++ {
		text := "test1"
		text2 := "test2"

		//Tell Go how many workers you want
		textWaitGroup.Add(2)

		//Worker functions just defined by putting "go" in front of them
		//Worker 1
		go func(message string) {
			defer textWaitGroup.Done()
			fmt.Println(message)
		}(text)

		//Worker 2
		go func(message2 string) {
			defer textWaitGroup.Done()
			fmt.Println(sum(1))
		}(text2)
	}
	textWaitGroup.Wait()
}
