package main

import (
	"strings"
	"fmt"
)

func shortenString(message string) func() string {
	return func() string {
		messageSlice := strings.Split(message, " ")
		wordLength := len(messageSlice)
		if wordLength < 1 {
			return "Nothing Left!"
		}
		messageSlice = messageSlice[:(wordLength - 1)]
		message = strings.Join(messageSlice, " ")
		return message
	}
}

func main() {
	myString := shortenString("Welcome to concurrency in Go! ...")
	fmt.Println(myString())
	fmt.Println(myString())
	fmt.Println(myString())
	fmt.Println(myString())
	fmt.Println(myString())
	fmt.Println(myString())
	fmt.Println(myString())
	fmt.Println(myString())
	fmt.Println(myString())
}
