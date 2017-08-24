package tester

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/rcliao/sql-unit-test/runner"
)

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
	Database string `json:"database"`
	Host     string `json:"host"`
}

// Run runs through a list of statements in the submission and compare to test cases
func Run(runner runner.Runner, statements, setupStatements, teardownStatements []Statement, testcases map[string][]map[string]string) (bool, error) {
	var pass = true
	for _, statement := range setupStatements {
		if err := runner.Execute(statement.Text); err != nil {
			return false, err
		}
	}
	for i, statement := range statements {
		result, err := runner.Query(statement.Text)
		if err != nil {
			return false, err
		}
		expected := testcases[strconv.Itoa(i+1)]
		if !reflect.DeepEqual(result, expected) {
			fmt.Printf("Test case %d: expected %v but got %v", i, expected, result)
			pass = false
		}
	}
	// teardown
	for _, statement := range teardownStatements {
		if err := runner.Execute(statement.Text); err != nil {
			return false, err
		}
	}
	return pass, nil
}
