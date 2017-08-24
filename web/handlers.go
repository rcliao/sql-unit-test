package web

import (
	"fmt"
	"html/template"
	"net/http"
)

// Hello says hello
func Hello() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, SQL-Unit-Test Server!")
	})
}

// Index renders the index page for submitting SQL queries to test
func Index() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./web/templates/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, nil)
	})
}
