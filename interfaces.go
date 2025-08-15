package main

import (
	"fmt"
)

// This is the main Notifer interface. Determines whether a variable has the function if it matches this type
type Notifier interface {
	SendNotification(message string)
}

type Multiplier interface {
	MultiplyDigits() int
}

// Structs to determine the object types
type EmailNotifier struct {
	recipientEmail string
}

type SMSNotifier struct {
	phoneNumber string
}

type PushNotifier struct {
	deviceID string
}

type numberMultiplier struct {
	x int
	y int
}

// Defining the methods on each of the structs. The first set of brackets needs to reference the struct not the interface
func (e EmailNotifier) SendNotification(message string) {
	fmt.Println("Email to", e.recipientEmail+":", message)
}

func (s SMSNotifier) SendNotification(message string) {
	fmt.Println("You have sent the text:", message)
}

func (p PushNotifier) SendNotification(message string) {
	fmt.Println("Push notification sent to device", p.deviceID+":", message)
}

func (m numberMultiplier) MultiplyDigits() int {
	return m.x * m.y
}

func main() {
	var emailNotification Notifier
	emailNotification = EmailNotifier{recipientEmail: "jordan"}
	emailNotification.SendNotification("Hiya")

	var SMSNotification Notifier
	SMSNotification = SMSNotifier{phoneNumber: "0414927621"}
	SMSNotification.SendNotification("This is a test push notification")

	var pushNotification Notifier
	pushNotification = PushNotifier{deviceID: "2813"}
	pushNotification.SendNotification("Test notification number 4")

	var myTwoNumbers Multiplier
	myTwoNumbers = numberMultiplier{x: 4, y: 5}
	fmt.Println(myTwoNumbers.MultiplyDigits())

	var yourTwoNumbers Multiplier
	yourTwoNumbers = numberMultiplier{x: 7, y: 9}
	fmt.Println(yourTwoNumbers.MultiplyDigits())
}
