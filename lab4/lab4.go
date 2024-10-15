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
	result int
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
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		show_error(w)
		return
	}
	var num1, num2 int
	num1, err1 := strconv.Atoi(parts[1])
	num2, err2 := strconv.Atoi(parts[2])
	validOps := map[string]bool{"add": true, "sub": true, "mul": true, "div": true, "gcd": true, "lcm": true}
	if !validOps[parts[0]] || err1 != nil || err2 != nil {
		show_error(w)
		return
	}
	fmt.Fprintf(w, "%d, %d", num1, num2)
	


}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8084", nil))
}
