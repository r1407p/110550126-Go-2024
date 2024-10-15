package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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



}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8084", nil))
}
