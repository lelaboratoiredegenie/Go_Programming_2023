package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Calculator(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// fmt.Printf("parts[0]: %s\n", parts[0])
	// fmt.Printf("parts[1]: %s\n", parts[1])
	if len(parts) != 4 {
		fmt.Fprintf(w, "Error!")
		return
	}

	operation := parts[1]
	num1Str := parts[2]
	num2Str := parts[3]

	num1, err := strconv.Atoi(num1Str)
	if err != nil {
		fmt.Fprintf(w, "Error!")
		return
	}
	num2, err := strconv.Atoi(num2Str)
	if err != nil {
		fmt.Fprintf(w, "Error!")
		return
	}

	var result int
	switch strings.ToLower(operation) {
	case "add":
		result = num1 + num2
		fmt.Fprintf(w, "%d + %d = %d", num1, num2, result)
		return
	case "sub":
		result = num1 - num2
		fmt.Fprintf(w, "%d - %d = %d", num1, num2, result)
		return
	case "mul":
		result = num1 * num2
		fmt.Fprintf(w, "%d * %d = %d", num1, num2, result)
		return
	case "div":
		if num2 == 0 {
			fmt.Fprintf(w, "Error!")
			return
		}
		result = num1 / num2
		fmt.Fprintf(w, "%d / %d = %d, reminder = %d", num1, num2, result, num1-(num2*result))
		return
	default:
		fmt.Fprintf(w, "Error!")
		return
	}
}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8083", nil))
}
