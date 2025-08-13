package main

import (
	"fmt"
	"os"
	"bufio"
	"os/exec"
	"sync"
)

func main() {

	filePath := "site.txt"

	file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Unable to read file")
			fmt.Println(err)
		}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var pingWG sync.WaitGroup

	for scanner.Scan() {
		address := scanner.Text()
		if address == "" {
			continue
		}
		//fmt.Println(line)

		pingWG.Add(1)
		go func(addr string) {
			defer pingWG.Done()
			output, err := exec.Command("cmd", "/C", "ping", addr).CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(output))
		}(address)
	}
	pingWG.Wait()
}