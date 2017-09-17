package tester

import (
	"database/sql"
	"log"
	"math/rand"
	"reflect"

	"github.com/pkg/errors"

	"github.com/rcliao/sql-unit-test/db"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var databaseNameSize = 32

// TestCase represents each test case used against the Table
type TestCase struct {
	Index    string
	Content  []map[string]string
	Question string
}

// TestResult wraps the result of tests
type TestResult struct {
	Expected TestCase
	Actual   db.Table
	Pass     bool
	Error    error
}

// Statement represents each statement student submit
type Statement struct {
	Comment string
	Index   int
	Text    string
}

// Config for database connections
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Database string `json:"database"`
}

// Run runs through a list of statements in the submission and compare to test cases
func Run(
	sqlDB *sql.DB,
	statements, setupStatements, teardownStatements []Statement,
	testcases []TestCase,
	selectedQuestions []string,
) ([]TestResult, error) {
	var testResult = []TestResult{}
	// use random 32 characters database
	randomDatabaseName := getRandomString()

	tx, err := sqlDB.Begin()
	if err != nil {
		return testResult, errors.Wrap(err, "Have issue starting transaction ")
	}
	if _, err := tx.Exec("CREATE DATABASE " + randomDatabaseName); err != nil {
		return testResult, errors.Wrap(err, "Have issue creating database "+randomDatabaseName)
	}
	if _, err := tx.Exec("USE " + randomDatabaseName); err != nil {
		return testResult, errors.Wrap(err, "Have issue use random database")
	}
	defer func() {
		if _, err := tx.Exec("DROP DATABASE " + randomDatabaseName); err != nil {
			log.Println("Have issue dropping database", err)
		}
		if err := tx.Commit(); err != nil {
			log.Println("Have issue commiting transaction", err)
		}
	}()

	// Setup
	for _, statement := range setupStatements {
		if _, err := tx.Exec(statement.Text); err != nil {
			return testResult, errors.Wrap(err, "Have issue running setup statements")
		}
	}
	i := 0

	for _, expected := range testcases {
		if !stringInSlice(expected.Index, selectedQuestions) && len(selectedQuestions) > 0 {
			continue
		}
		if i >= len(statements) {
			testResult = append(
				testResult,
				TestResult{
					Expected: expected,
					Actual:   db.Table{},
					Pass:     false,
				},
			)
			i++
			continue
		}
		statement := statements[i]
		// TODO: detect if statement is a query or update
		table, err := db.Query(tx, statement.Text)
		// Query has syntax error
		if err != nil {
			testResult = append(
				testResult,
				TestResult{
					Expected: expected,
					Actual:   table,
					Pass:     false,
					Error:    err,
				},
			)
			i++
			continue
		}
		if !reflect.DeepEqual(table.Content, expected.Content) {
			testResult = append(
				testResult,
				TestResult{
					Expected: expected,
					Actual:   table,
					Pass:     false,
				},
			)
		} else {
			testResult = append(
				testResult,
				TestResult{
					Expected: expected,
					Actual:   table,
					Pass:     true,
				},
			)
		}
		i++
	}

	// teardown
	for _, statement := range teardownStatements {
		if _, err := tx.Exec(statement.Text); err != nil {
			return testResult, errors.Wrap(err, "Have issue running teardown statements")
		}
	}

	return testResult, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getRandomString() string {
	b := make([]rune, databaseNameSize)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
