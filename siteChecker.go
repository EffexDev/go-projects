package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

func worker(id int, jobs <-chan string, results chan<- string) {
	for job := range jobs {

		result, err := exec.Command("cmd", "/C", "ping", job).CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}

		//Simulate some time as can't ping using CMD from online Go compiler
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)

		results <- fmt.Sprintf("%s", result)
	}
}

func main() {
	jobs := make(chan string, 10)
	results := make(chan string, 10)

	for id := 1; id <= 3; id++ {
		go worker(id, jobs, results)
	}

	urls := []string{
		"www.google.com",
		"www.github.com",
		"www.facebook.com",
		"www.x.com",
		"www.buddytelco.com.au",
		"www.aussiebroadband.com.au",
		"www.telstra.com",
		"www.optus.com.au",
		"8.8.8.8",
	}

	for _, u := range urls {
		jobs <- u
	}
	close(jobs)

	for i := 0; i < len(urls); i++ {
		fmt.Println(<-results)
	}
}
