package main

import "fmt"

func main() {
	var n int64

	fmt.Print("Enter a number: ")
	fmt.Scanln(&n)

	result := Sum(n)
	fmt.Println(result)
}

func Sum(n int64) string {
	var result string
	var sum int64
	for i := int64(1); i <= n; i++ {
		if i % 7 == 0 {
			continue
		}
		if i == 1 {
			result = fmt.Sprintf("%d", i)
		} else {
			result = fmt.Sprintf("%s+%d", result, i)
		}
		sum += i
	}
	result = fmt.Sprintf("%s=%d", result, sum)
	return result
}