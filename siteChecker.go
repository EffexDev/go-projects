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
		}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var pingWG sync.WaitGroup

	for scanner.Scan() {
		address := scanner.Text()
		if address =="" {
			continue
		}
		//fmt.Println(line)

		pintWG.Add(1)
		go func(addr string) {
			defer pingWG.Done()
			cmd := exec.Command("cmd", "/C", "ping", addr).CombinedOutput()
			fmt.Println(string(output))
		}(address)
		
		
		output, _ := cmd
		fmt.Println(string(output))
	}
	pingWG.Wait()
}

