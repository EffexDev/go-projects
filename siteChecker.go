package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(id int, jobs <-chan string, results chan<- string) {
	for job := range jobs {

		// result, err := exec.Command("cmd", "/C", "ping", job).CombinedOutput()
		// if err != nil {
		// 	fmt.Println(err)
		// }

		//Simulate some time as can't ping using CMD from online Go compiler
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)

		results <- fmt.Sprintf("URL %s checked by worker %d", job, id)
	}
}

func main() {
	jobs := make(chan string)
	results := make(chan string)

	for id := 1; id <= 5; id++ {
		go worker(id, jobs, results)
	}

	urls := []string{
		"www.google.com",
		"www.github.com",
		"www.facebook.com",
		"www.x.com",
	}

	go func() {
		for _, u := range urls {
			jobs <- u
		}
		close(jobs)
	}()

	for i := 0; i < len(urls); i++ {
		fmt.Println(<-results)
	}
}
