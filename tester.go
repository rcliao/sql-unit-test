package tester

import (
	"database/sql"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/rcliao/sql-unit-test/db"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var databaseNameSize = 32

// DAO for database access object to define common behavior
type DAO interface {
	Query(query string) (Result, error)
	Exec(query string) error
}

// Result represents what returns from DAO
type Result struct {
	Query   string
	Content []map[string]string
}

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
	statements, setupStatements, teardownStatements, solutions []Statement,
	selectedQuestions []string,
) ([]TestResult, error) {
	var testResult = []TestResult{}

	i := 0
	tables, errs, err := ExecuteStatements(sqlDB, setupStatements, teardownStatements, statements)
	solutionTables, _, err := ExecuteStatements(sqlDB, setupStatements, teardownStatements, solutions)
	testcases := ConvertTablesToTestCases(solutionTables)

	if err != nil {
		return testResult, err
	}

	for _, expected := range testcases {
		if !stringInSlice(expected.Index, selectedQuestions) && len(selectedQuestions) > 0 {
			continue
		}
		if i >= len(tables) {
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
		table := tables[i]
		err := errs[i]
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

	return testResult, nil
}

// ExecuteStatements takes arguments to build a list of tables
func ExecuteStatements(sqlDB *sql.DB, setupStatements, teardownStatements, statements []Statement) ([]db.Table, []error, error) {
	result := []db.Table{}
	errs := []error{}

	// use random 32 characters database
	randomDatabaseName := getRandomString()

	tx, err := sqlDB.Begin()
	if err != nil {
		return result, errs, errors.Wrap(err, "Have issue starting transaction ")
	}

	if _, err := tx.Exec("CREATE DATABASE " + randomDatabaseName); err != nil {
		return result, errs, errors.Wrap(err, "Have issue creating database "+randomDatabaseName)
	}
	if _, err := tx.Exec("USE " + randomDatabaseName); err != nil {
		return result, errs, errors.Wrap(err, "Have issue use random database")
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
			return result, errs, errors.Wrap(err, "Have issue running setup statements")
		}
	}

	var errAccumulator error

	// execute all submitted statements and store them under tables
	for _, statement := range statements {
		// IDEA: it's probably better to create a list of READ query for reading check
		if strings.Index(strings.ToLower(statement.Text), "select") == 0 || strings.Index(strings.ToLower(statement.Text), "describe") == 0 {
			table, err := db.Query(tx, statement.Text)
			result = append(result, table)
			if errAccumulator != nil && err != nil {
				err = errors.Wrap(err, errAccumulator.Error())
			}
			if errAccumulator != nil && err == nil {
				err = errAccumulator
			}
			errs = append(errs, err)
			errAccumulator = nil
			continue
		}
		_, errAccumulator = tx.Exec(statement.Text)
		if errAccumulator != nil {
			errAccumulator = errors.Wrap(errAccumulator, "Query \""+statement.Text+"\" has syntax error.")
		}
	}

	// teardown
	for _, statement := range teardownStatements {
		if _, err := tx.Exec(statement.Text); err != nil {
			return result, errs, errors.Wrap(err, "Have issue running teardown statements")
		}
	}

	return result, errs, nil
}

// ConvertTablesToTestCases takes table content and convert it into TestCases
func ConvertTablesToTestCases(tables []db.Table) []TestCase {
	testcases := []TestCase{}
	for i, t := range tables {
		testcases = append(testcases, TestCase{
			Index:    strconv.Itoa(i + 1),
			Content:  t.Content,
			Question: "",
		})
	}

	return testcases
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
