package tester

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"

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
func Run(sqlDB *sql.DB, statements, setupStatements, teardownStatements []Statement, testcases []TestCase) ([]TestResult, error) {
	var testResult = []TestResult{}
	// use UUID to create temp database
	randomDatabaseName := getRandomString()
	if _, err := sqlDB.Exec("CREATE DATABASE " + randomDatabaseName); err != nil {
		fmt.Println(err)
		return testResult, err
	}
	if _, err := sqlDB.Exec("USE " + randomDatabaseName); err != nil {
		fmt.Println(err)
		return testResult, err
	}
	defer func() {
		if _, err := sqlDB.Exec("DROP DATABASE " + randomDatabaseName); err != nil {
			fmt.Println(err)
		}
	}()

	// Setup
	for _, statement := range setupStatements {
		if _, err := sqlDB.Exec(statement.Text); err != nil {
			return testResult, err
		}
	}

	for i, expected := range testcases {
		if i >= len(statements) {
			testResult = append(
				testResult,
				TestResult{
					Expected: expected,
					Actual:   db.Table{},
					Pass:     false,
				},
			)
			continue
		}
		statement := statements[i]
		// TODO: detect if statement is a query or update
		table, err := db.Query(sqlDB, statement.Text)
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
	}

	// teardown
	for _, statement := range teardownStatements {
		if _, err := sqlDB.Exec(statement.Text); err != nil {
			return testResult, err
		}
	}

	return testResult, nil
}

func getRandomString() string {
	b := make([]rune, databaseNameSize)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
