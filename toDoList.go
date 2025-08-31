package main

import (
	"fmt"
	"slices"
	"strconv"
	"os"
	"os/exec"
	"runtime"
	"bufio"
	"strings"
)

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

type toDoItem struct {
	number     int
	details    string
	completion bool
}

type toDoList []toDoItem

func ClearTerminal() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default: // Linux, macOS, etc.
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (t *toDoList) AddItem(taskDetails string) {
	newItem := toDoItem{
		number:  len(*t) + 1,
		details: taskDetails,
	}
	*t = append(*t, newItem)
}

func (t *toDoList) DeleteItem(taskNumber int) {
	*t = slices.Delete(*t, taskNumber-1, taskNumber)
}

func (t *toDoList) CompleteItem(taskNumber int) {
	(*t)[taskNumber-1].completion = !(*t)[taskNumber-1].completion
}

func main() {
	ClearTerminal()
	todos := toDoList{}

	var menuChoice string
	var taskToAdd string
	var itemToComplete string

	for {
		fmt.Println("\nWelcome to your daily task list!\n\nPlease select an option:")
		fmt.Println(ColorBlue + "1. Add new task" + ColorReset)
		fmt.Println(ColorYellow + "2. List all tasks" + ColorReset)
		fmt.Println(ColorGreen + "3. Complete a task" + ColorReset)

		fmt.Scan(&menuChoice)
		if menuChoice == "" {
			fmt.Println("Please enter a number")
		}

		choice, err := strconv.Atoi(menuChoice)
		if err != nil {
			fmt.Println("Failed to parse choice to integer")
		}

		switch choice {
		case 1:
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
			fmt.Println("Enter the task you would like to add.")
			taskToAdd, _ = reader.ReadString('\n')
			taskToAdd = strings.TrimSpace(taskToAdd)
			todos.AddItem(taskToAdd)
			ClearTerminal()
			fmt.Println("\n---------\nTask added. Current list:\n---------\n")
			for _, t := range todos {
				if t.completion == true {
					fmt.Printf("%d. %s ", t.number, t.details)
					fmt.Printf(ColorGreen + "- Complete\n" + ColorReset)
				} else {
					fmt.Printf("%d. %s ", t.number, t.details)
					fmt.Printf(ColorRed + "- Incomplete\n" + ColorReset)
				}
			}
		case 2:
			ClearTerminal()
			fmt.Println("\n---------\nCurrent task list:\n---------\n")
			for _, t := range todos {
				if t.completion == true {
					fmt.Printf("%d. %s ", t.number, t.details)
					fmt.Printf(ColorGreen + "- Complete\n" + ColorReset)
				} else {
					fmt.Printf("%d. %s ", t.number, t.details)
					fmt.Printf(ColorRed + "- Incomplete\n" + ColorReset)
				}
			}
		case 3:
			fmt.Println("\n---------\nWhich task would you like to complete? (Enter the task number)\n---------\n")
			fmt.Scan(&itemToComplete)

			num, err := strconv.Atoi(itemToComplete)
			if err != nil {
				fmt.Println(err)
			}

			if num >= len(todos)+1 {
				ClearTerminal()
				fmt.Println("Task does not exist")
				continue
			}

			todos.CompleteItem(num)
			ClearTerminal()
			fmt.Println("---------")
			fmt.Println("Task Completed. Current tasks:")
			for _, t := range todos {
				if t.completion == true {
					fmt.Printf("%d. %s ", t.number, t.details)
					fmt.Printf(ColorGreen + "- Complete\n" + ColorReset)
				} else {
					fmt.Printf("%d. %s ", t.number, t.details)
					fmt.Printf(ColorRed + "- Incomplete\n" + ColorReset)
				}
			}
		}
	}
}
