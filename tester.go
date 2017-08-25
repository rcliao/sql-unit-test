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
func Run(runner runner.Runner, statements, setupStatements, teardownStatements []Statement, testcases map[string][]map[string]string) ([]string, error) {
	var failedTestCases = []string{}
	for _, statement := range setupStatements {
		if err := runner.Execute(statement.Text); err != nil {
			return failedTestCases, err
		}
	}
	for i, statement := range statements {
		result, err := runner.Query(statement.Text)
		if err != nil {
			return failedTestCases, err
		}
		expected := testcases[strconv.Itoa(i+1)]
		if !reflect.DeepEqual(result, expected) {
			failedTestCases = append(
				failedTestCases,
				fmt.Sprintf(
					"Failed test case %d: expected %v but got %v\n",
					i,
					expected,
					result,
				),
			)
		}
	}
	// teardown
	for _, statement := range teardownStatements {
		if err := runner.Execute(statement.Text); err != nil {
			return failedTestCases, err
		}
	}
	return failedTestCases, nil
}
