package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {

	filePath := "site.txt"

	file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Unable to read file")
		}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

}

