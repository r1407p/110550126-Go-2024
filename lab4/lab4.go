package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// TODO: Create a struct to hold the data sent to the template
type Data struct {
	op string
	num1 int
	num2 int
	Result int
	Expression string
}

func show_error(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("error.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func gcd(a int, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a % b)
}

func lcm(a int, b int) int {
	return a * b / gcd(a, b)
}


func get_result(op string, num1 int, num2 int) Data {
	var result Data
	result.op = op
	result.num1 = num1
	result.num2 = num2

	switch op {
		case "add":
			result.Result = num1 + num2
			result.Expression = fmt.Sprintf("%d + %d", num1, num2)
		case "sub":
			result.Result = num1 - num2
			result.Expression = fmt.Sprintf("%d - %d", num1, num2)
		case "mul":
			result.Result = num1 * num2
			result.Expression = fmt.Sprintf("%d * %d", num1, num2)
		case "div":
			result.Result = num1 / num2
			// remainder := num1 % num2
			result.Expression = fmt.Sprintf("%d / %d", num1, num2)
		case "gcd":
			result.Result = gcd(num1, num2)
			result.Expression = fmt.Sprintf("GCD(%d, %d)", num1, num2)
		case "lcm":
			result.Result = lcm(num1, num2)
			result.Expression = fmt.Sprintf("LCM(%d, %d)", num1, num2)
	}
	return result

}
func Calculator(w http.ResponseWriter, r *http.Request) {
	// TODO: Finish this function
	op := r.URL.Query().Get("op")
	num1Str := r.URL.Query().Get("num1")
	num2Str := r.URL.Query().Get("num2")

	if op == "" || num1Str == "" || num2Str == "" {
		show_error(w)
		return
	}
	var num1, num2 int
	num1, err1 := strconv.Atoi(num1Str)
	num2, err2 := strconv.Atoi(num2Str)

	validOps := map[string]bool{"add": true, "sub": true, "mul": true, "div": true, "gcd": true, "lcm": true}
	if !validOps[op] || err1 != nil || err2 != nil {
		show_error(w)
		return
	}

	if op == "div" && num2 == 0 {
		show_error(w)
		return
	}
	var result Data
	result = get_result(op, num1, num2)

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, result)

}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8084", nil))
}
