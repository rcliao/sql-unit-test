package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	tester "github.com/rcliao/sql-unit-test"
	"github.com/rcliao/sql-unit-test/parser"
)

var subjectFolder = "./subjects"

// Page is simple DTO to transfer multiple information to index.html
type Page struct {
	Instruction template.HTML
	Subject     string
	Testcases   []tester.TestCase
}

// SummaryPage is simple DTO to transfer info to result.html
type SummaryPage struct {
	Results        []tester.TestResult
	NumberOfPasses int
}

// Hello says hello
func Hello() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, SQL-Unit-Test Server!")
	})
}

// HealthCheck returns the healthcheck for the critical resources
func HealthCheck(sqlDB *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := sqlDB.Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Fprintln(w, "Healthy")
	})
}

// Static serves the static assets (js & css)
func Static() http.Handler {
	return http.StripPrefix("/static", http.FileServer(http.Dir("./web/static")))
}

// Index renders the index page for submitting SQL queries to test
func Index() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subject := vars["subject"]
		testcases := []tester.TestCase{}
		instruction := template.HTML("")

		if subject != "" {
			content, err := ioutil.ReadFile(subjectFolder + "/" + subject + "/instruction.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			instruction = template.HTML(content)
			testcasesContent, err := ioutil.ReadFile(subjectFolder + "/" + subject + "/testcase.json")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			testcases, err = parser.ParseTestCases(string(testcasesContent))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		t, err := template.ParseFiles("./web/templates/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, Page{
			Instruction: instruction,
			Subject:     subject,
			Testcases:   testcases,
		})
	})
}

// RunTest handles the query from the request param to test them
func RunTest(sqlDB *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		submission := r.FormValue("statements")
		subject := r.FormValue("subject")
		selectedQuestions := r.Form["question_numbers"]

		defer func() {
			log.Println("Executed SQL Unit Test", r)
		}()

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
		testResult, err := tester.Run(
			sqlDB,
			statements,
			setupStatements,
			teardownStatements,
			testCases,
			selectedQuestions,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error while running test cases", err)
			return
		}
		numOfPasses := 0
		for _, t := range testResult {
			if t.Pass {
				numOfPasses++
			}
		}

		t, err := template.ParseFiles("./web/templates/result.html")
		summaryDTO := SummaryPage{
			Results:        testResult,
			NumberOfPasses: numOfPasses,
		}
		if err != nil {
			panic(err)
		}
		t.Execute(w, summaryDTO)
	})
}
