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
    Expression string
    Result  string
}

func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func Calculator(w http.ResponseWriter, r *http.Request) {
	// TODO: Finish this function
	num1Str := r.URL.Query().Get("num1")
	num2Str := r.URL.Query().Get("num2")
	operation := r.URL.Query().Get("op")

	num1, err := strconv.Atoi(num1Str)
	if err != nil {
		http.ServeFile(w, r, "error.html")
		return
	}
	num2, err := strconv.Atoi(num2Str)
	if err != nil {
		http.ServeFile(w, r, "error.html")
		return
	}

	var result int
	var e string
	var s string
	switch (operation) {
	case "add":
		result = num1 + num2

		e = fmt.Sprintf("%d + %d", num1, num2)
		s = fmt.Sprintf("%d", result)

		data := Data{
			Expression: e,
    		Result: s,
		}

		err := template.Must(template.ParseFiles("index.html")).Execute(w, data)

    	if err != nil {
    	    http.ServeFile(w, r, "error.html")
    	    return
    	}

		return
	case "sub":
		result = num1 - num2
		
		e = fmt.Sprintf("%d - %d", num1, num2)
		s = fmt.Sprintf("%d", result)

		data := Data{
			Expression: e,
    		Result: s,
		}

		err := template.Must(template.ParseFiles("index.html")).Execute(w, data)

    	if err != nil {
    	    http.ServeFile(w, r, "error.html")
    	    return
    	}

		return
	case "mul":
		result = num1 * num2
		e = fmt.Sprintf("%d * %d", num1, num2)
		s = fmt.Sprintf("%d", result)

		data := Data{
			Expression: e,
    		Result: s,
		}

		err := template.Must(template.ParseFiles("index.html")).Execute(w, data)

    	if err != nil {
    	    http.ServeFile(w, r, "error.html")
    	    return
    	}
		return
	case "div":
		if num2 == 0 {
			http.ServeFile(w, r, "error.html")
			return
		}
		result = num1 / num2
		e = fmt.Sprintf("%d / %d", num1, num2)
		s = fmt.Sprintf("%d", result)

		data := Data{
			Expression: e,
    		Result: s,
		}

		err := template.Must(template.ParseFiles("index.html")).Execute(w, data)

    	if err != nil {
    	    http.ServeFile(w, r, "error.html")
    	    return
    	}
		return

	case "lcm":
		result = (num1 * num2) / gcd(num1, num2)

		e = fmt.Sprintf("LCM(%d, %d)", num1, num2)
		s = fmt.Sprintf("%d", result)

		data := Data{
			Expression: e,
    		Result: s,
		}

		err := template.Must(template.ParseFiles("index.html")).Execute(w, data)

    	if err != nil {
    	    http.ServeFile(w, r, "error.html")
    	    return
    	}

		return
	case "gcd":
		result = gcd(num1, num2)

		e = fmt.Sprintf("GCD(%d, %d)", num1, num2)
		s = fmt.Sprintf("%d", result)

		data := Data{
			Expression: e,
    		Result: s,
		}

		err := template.Must(template.ParseFiles("index.html")).Execute(w, data)

    	if err != nil {
    	    http.ServeFile(w, r, "error.html")
    	    return
    	}

		return
	default:
		http.ServeFile(w, r, "error.html")
		return
	}
}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8084", nil))
}
