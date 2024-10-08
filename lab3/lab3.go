package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Calculator(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		fmt.Fprintf(w, "Error!")
		return
	}
	var num1, num2 int
	// num1, err1 := strconv.ParseFloat(parts[1], 64)
	// num2, err2 := strconv.ParseFloat(parts[2], 64)
	num1, err1 := strconv.Atoi(parts[1])
	num2, err2 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil {
		fmt.Fprintf(w, "Error!")
		return
	}

	var result int
	switch parts[0] {
	case "add":
		result = num1 + num2
		fmt.Fprintf(w, "%d + %d = %d", num1, num2, result)
	case "sub":
		result = num1 - num2
		fmt.Fprintf(w, "%d - %d = %d", num1, num2, result)
	case "mul":
		result = num1 * num2
		fmt.Fprintf(w, "%d * %d = %d", num1, num2, result)
	case "div":
		if num2 == 0 {
			fmt.Fprintf(w, "Error!")
			return
		}
		result = num1 / num2
		remainder := num1 % num2
		fmt.Fprintf(w, "%d / %d = %d, reminder = %d", num1, num2, result, remainder)
	default:
		fmt.Fprintf(w, "Error!")
		return
	}

	


}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8083", nil))
}