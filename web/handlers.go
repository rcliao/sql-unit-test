package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	tester "github.com/rcliao/sql-unit-test"
	"github.com/rcliao/sql-unit-test/db"
	"github.com/rcliao/sql-unit-test/parser"
)

var subjectFolder = "./subjects"
var solutionCache = make(map[string][]tester.TestCase)

// TODO: convert the following to be more like different setup
var subjectTypes = map[string]string{
	"homework-3": "mongo",
	"homework-4": "mongo",
}

func getSubjectType(subject string) string {
	t, okay := subjectTypes[subject]
	if !okay {
		return "sql"
	}
	return t
}

// Page is simple DTO to transfer multiple information to index.html
type Page struct {
	Instruction template.HTML
	Subject     string
	TestCases   []tester.TestCase
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
func Index(factory *db.Factory) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subject := vars["subject"]
		if subject == "" {
			subject = "exercise-1"
		}
		dao := factory.CreateDAO(getSubjectType(subject))
		testcases := []tester.TestCase{}
		instruction := template.HTML("")

		if subject != "" {
			content, err := ioutil.ReadFile(subjectFolder + "/" + subject + "/instruction.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			instruction = template.HTML(content)
			solutionContent, err := ioutil.ReadFile(subjectFolder + "/" + subject + "/solution.sql")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			solutions := parser.ParseSQL(string(solutionContent), "#")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if cache, okay := solutionCache[subject]; !okay {
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
				tables, _, err := dao.ExecuteStatements(setupStatements, teardownStatements, solutions)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(tables)
				testcases = tester.ConvertTablesToTestCases(tables)
				solutionCache[subject] = testcases
			} else {
				testcases = cache
			}
		}

		t, err := template.ParseFiles("./web/templates/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, Page{
			Instruction: instruction,
			Subject:     subject,
			TestCases:   testcases,
		})
	})
}

// RunTest handles the query from the request param to test them
func RunTest(factory *db.Factory) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		submission := r.FormValue("statements")
		subject := r.FormValue("subject")
		dao := factory.CreateDAO(getSubjectType(subject))
		selectedQuestions := r.Form["question_numbers"]

		defer func() {
			log.Println("Executed SQL Unit Test", r)
		}()

		// parsing content
		solutionContent, err := ioutil.ReadFile(subjectFolder + "/" + subject + "/solution.sql")
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
		solutions := parser.ParseSQL(string(solutionContent), "#")
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
		// hacking around the setup for now
		if getSubjectType(subject) == "mongo" {
			for i, statement := range setupStatements {
				setupStatements[i].Text = strings.Replace(statement.Text, "{filePath}", subjectFolder+"/"+subject+"", -1)
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
			dao,
			statements,
			setupStatements,
			teardownStatements,
			solutions,
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
