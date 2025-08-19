package main

import (
	"fmt"
	"slices"
	"strconv"
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
	todos := toDoList{
		{1, "Create a OOP todo in Go", false},
	}

	menuChoice := "3"
	itemToDelete := "1"
	taskToAdd := "Publish to Github"

	fmt.Println("Welcome to your daily task list!\n\nPlease select an option:")
	fmt.Println(ColorBlue + "1. Add new task" + ColorReset)
	fmt.Println(ColorYellow + "2. List all tasks" + ColorReset)
	fmt.Println(ColorGreen + "3. Complete a task" + ColorReset)
	fmt.Println(ColorRed + "4. Delete a task" + ColorReset)

	// fmt.Scan(&menuChoice)
	// if menuChoice == "" {
	// 	fmt.Println("Please enter a number")
	// }

	choice, err := strconv.Atoi(menuChoice)
	if err != nil {
		fmt.Println("Failed to parse choice to integer")
		return
	}

	switch choice {
	case 1:
		todos.AddItem(taskToAdd)
		fmt.Println("\n---------\nTask added. Current list:\n---------\n")
		for _, t := range todos {
			fmt.Printf("%d. %s - Completed: %v\n", t.number, t.details, t.completion)
		}
	case 2:
		fmt.Println("\n---------\nCurrent task list:\n---------\n")
		for _, t := range todos {
			fmt.Printf("%d. %s - Completed: %v\n", t.number, t.details, t.completion)
		}
	case 3:
		itemToComplete := "1"
		fmt.Println("\n---------\nWhich task would you like to complete? (Enter the task number)\n---------\n")
		fmt.Scan(&itemToComplete)

		num, err := strconv.Atoi(itemToComplete)
		if err != nil {
			fmt.Println(err)
		}
		todos.CompleteItem(num)
		fmt.Println("---------")
		fmt.Println("Task Completed. Current tasks:")
		for _, t := range todos {
			fmt.Printf("%d. %s - Completed: %v\n", t.number, t.details, t.completion)
		}
	case 4:
		num, err := strconv.Atoi(itemToDelete)
		if err != nil {
			fmt.Println(err)
		}
		todos.DeleteItem(num)
		fmt.Println(todos)
	}

}
