package main

import (
	"fmt"
	"math/big"
	"strconv"
	"syscall/js"
)

func CheckPrime(this js.Value, args []js.Value) interface{} {
	// TODO: Check if the number is prime
	input_number := js.Global().Get("document").Call("getElementById", "value").Get("value").String()
	fmt.Println("Input number: ", input_number)
	number, _ := strconv.Atoi(input_number)
	big_number := big.NewInt(int64(number))
	is_prime := big_number.ProbablyPrime(0)
	if is_prime {
		js.Global().Get("document").Call("getElementById", "answer").Set("innerHTML", "It's prime")
	} else {
		js.Global().Get("document").Call("getElementById", "answer").Set("innerHTML", "It's not prime")
	}
	return nil
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