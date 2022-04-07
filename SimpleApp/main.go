package main

import (
	"html/template"
	"net/http"
	"strconv"

	"./computation"
)

type ViewData struct {
	Answer string
}

func main() {

	http.HandleFunc("/", homepage)
	http.ListenAndServe(":80", nil)

}

func homepage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		data := ViewData{Answer: `Решение: Введите коэффициенты и нажмите "посчитать"`}
		tmpl, _ := template.ParseFiles("static/index.html")
		tmpl.Execute(w, data)

	}

	if r.Method == "POST" {

		a, _ := strconv.ParseFloat(r.FormValue("a"), 64)
		b, _ := strconv.ParseFloat(r.FormValue("b"), 64)
		c, _ := strconv.ParseFloat(r.FormValue("c"), 64)
		answer := computation.Solve(a, b, c)
		data := ViewData{Answer: answer}
		tmpl, _ := template.ParseFiles("static/index.html")
		tmpl.Execute(w, data)

	}

}
