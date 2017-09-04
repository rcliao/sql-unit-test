package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	tester "github.com/rcliao/sql-unit-test"
	"github.com/rcliao/sql-unit-test/parser"
)

var subjectFolder = "./subjects"

// Hello says hello
func Hello() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, SQL-Unit-Test Server!")
	})
}

// Static serves the static assets (js & css)
func Static() http.Handler {
	return http.StripPrefix("/static", http.FileServer(http.Dir("./web/static")))
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

// RunTest handles the query from the request param to test them
func RunTest(sqlDB *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		submission := r.FormValue("statements")
		subject := r.FormValue("subject")

		// parsing content
		testCasesContent, err := ioutil.ReadFile(subjectFolder + "/" + subject + "/testcase.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		setupContent, err := ioutil.ReadFile(subjectFolder + "/" + subject + "/setup.sql")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		teardownContent, err := ioutil.ReadFile(subjectFolder + "/" + subject + "/teardown.sql")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		testCases, err := parser.ParseTestCases(string(testCasesContent))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var setupStatements = []tester.Statement{}
		if string(setupContent) != "" {
			setupStatements = parser.ParseSQL(string(setupContent), "#")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		var teardownStatements = []tester.Statement{}
		if string(teardownContent) != "" {
			teardownStatements = parser.ParseSQL(string(teardownContent), "#")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		statements := parser.ParseSQL(submission, "#")
		testResult, err := tester.Run(sqlDB, statements, setupStatements, teardownStatements, testCases)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error while running test cases", err)
			return
		}

		t, err := template.ParseFiles("./web/templates/result.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, testResult)
	})
}
