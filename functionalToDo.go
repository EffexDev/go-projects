package main

import (
	"fmt"
	"slices"
	"strconv"
)

/*
   Things to implement:
       Creating a to-do item - Completed
       Deleting a to-do item - Completed
       Marking an item completed - Completed
       Listing all items - Completed
*/

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

type todoItem struct {
	itemNumber  int
	taskDetails string
	completed   bool
}

func createItem(todoList []todoItem, newTaskDetails string, completion bool) []todoItem {
	newItem := todoItem{
		itemNumber:  len(todoList) + 1,
		taskDetails: newTaskDetails,
		completed:   completion,
	}
	return append(todoList, newItem)
}

func completeTask(todoList []todoItem, taskNumber int) []todoItem {
	todoList[taskNumber-1].completed = !todoList[taskNumber-1].completed
	return todoList
}

func deleteTask(todoList []todoItem, taskNumber int) []todoItem {
	return slices.Delete(todoList, taskNumber-1, taskNumber)
}

func main() {
	todoList := []todoItem{
		{itemNumber: 1, taskDetails: "Make a to-do in Go", completed: false},
		{itemNumber: 2, taskDetails: "Publish to github", completed: false},
	}

	//test data as fmt.Scan() doesn't work in playground
	menuChoice := "4"
	newTaskDetails := "Make profit"
	newTaskCompletion := false

	fmt.Println("Welcome to your daily task list!\n\nPlease select an option:")
	fmt.Println(ColorBlue + "1. Add new task" + ColorReset)
	fmt.Println(ColorYellow + "2. List all tasks" + ColorReset)
	fmt.Println(ColorGreen + "3. Complete a task" + ColorReset)
	fmt.Println(ColorRed + "4. Delete a task" + ColorReset)

	num, err := strconv.Atoi(menuChoice)
	if err != nil {
		fmt.Println("Failed to parse choice to integer")
		return
	}

	// fmt.Scan(&menuChoice)
	// if menuChoice == "" {
	// 	fmt.Println("Please enter a number")
	// }

	switch num {
	case 1:
		todoList = createItem(todoList, newTaskDetails, newTaskCompletion)
		fmt.Println("----------")
		fmt.Println("Task added successfully. Current tasks:")
		for _, t := range todoList {
			fmt.Printf("%d. %s - Completed: %v\n", t.itemNumber, t.taskDetails, t.completed)
		}
	case 2:
		fmt.Println("----------")
		fmt.Println("Current tasks to complete:")
		for _, t := range todoList {
			fmt.Printf("%d. %s - Completed: %v\n", t.itemNumber, t.taskDetails, t.completed)
		}
	case 3:
		itemToComplete := "1"
		fmt.Println("----------")
		fmt.Println("Which task would you like to complete? (Enter the task number)")
		fmt.Scan(&itemToComplete)

		num, err := strconv.Atoi(itemToComplete)
		if err != nil {
			fmt.Println(err)
		}
		completeTask(todoList, num)
		fmt.Println("----------")
		fmt.Println("Task Completed. Current tasks:")
		for _, t := range todoList {
			fmt.Printf("%d. %s - Completed: %v\n", t.itemNumber, t.taskDetails, t.completed)
		}
	case 4:
		itemToDelete := "1"
		fmt.Println("----------")
		fmt.Println("Which task would you like to delete? (Enter the task number)")
		// fmt.Scan(&itemToDelete)

		num, err := strconv.Atoi(itemToDelete)
		if err != nil {
			fmt.Println(err)
			return
		}
		todoList = deleteTask(todoList, num)
		fmt.Println("----------")
		fmt.Println("Task Deleted. Current tasks to complete:")
		for _, t := range todoList {
			fmt.Printf("%d. %s - Completed: %v\n", t.itemNumber, t.taskDetails, t.completed)
		}
	}
}
