package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	tester "github.com/rcliao/sql-unit-test"
	"github.com/rcliao/sql-unit-test/parser"
	"github.com/rcliao/sql-unit-test/runner"
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

// RunTest handles the query from the request param to test them
func RunTest() http.HandlerFunc {
	testCasesContent, err := ioutil.ReadFile("./testcase.json")
	if err != nil {
		panic(err)
	}
	setupContent, err := ioutil.ReadFile("./setup.sql")
	if err != nil {
		panic(err)
	}
	teardownContent, err := ioutil.ReadFile("./teardown.sql")
	if err != nil {
		panic(err)
	}
	testCases, err := parser.ParseTestCases(string(testCasesContent))
	if err != nil {
		panic(err)
	}
	var setupStatements = []tester.Statement{}
	if string(setupContent) != "" {
		setupStatements = parser.ParseSQL(string(setupContent), "#")
		if err != nil {
			panic(err)
		}
	}
	var teardownStatements = []tester.Statement{}
	if string(teardownContent) != "" {
		teardownStatements = parser.ParseSQL(string(teardownContent), "#")
		if err != nil {
			panic(err)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := getDB()
		runner := runner.NewMySQLRunner(db)

		submission := r.FormValue("statements")
		statements := parser.ParseSQL(string(submission), "#")
		failedTestCases, err := tester.Run(runner, statements, setupStatements, teardownStatements, testCases)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var result = "All test passes!"
		if len(failedTestCases) > 0 {
			result = strings.Join(failedTestCases, "\n")
		}

		fmt.Fprintln(w, result)
	})
}

func getDB() *sql.DB {
	// TODO: generate a random DB string
	defaultProtocol := "tcp"
	defaultPort := "3306"
	sqlDSN := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		defaultProtocol,
		os.Getenv("MYSQL_HOST"),
		defaultPort,
		os.Getenv("MYSQL_DB"),
	)

	db, err := sql.Open("mysql", sqlDSN)
	if err != nil {
		panic(err)
	}

	return db
}
