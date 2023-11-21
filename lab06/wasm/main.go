package main

import (
	"fmt"
	"math/big"
	"strconv"
	"syscall/js"
)

func CheckPrime(this js.Value, args []js.Value) interface{} {
	// TODO: Check if the number is prime
	input := js.Global().Get("value").Get("value").String()
	inputValue, _ := strconv.ParseInt(input, 10, 64)
	num := big.NewInt(inputValue) // Replace 17 with any number you want to test

	isProbablyPrime := num.ProbablyPrime(0)

	fmt.Println(isProbablyPrime)

	if isProbablyPrime {
		js.Global().Get("answer").Set("innerText", "It's prime")
	} else {
		js.Global().Get("answer").Set("innerText", "It's not prime")
	}

	return "hi"
}

func registerCallbacks() {
	// TODO: Register the function CheckPrime
	js.Global().Set("CheckPrime", js.FuncOf(CheckPrime))
}

func main() {
	fmt.Println("Golang main function executed")
	registerCallbacks()

	//need block the main thread forever
	select {}
}
