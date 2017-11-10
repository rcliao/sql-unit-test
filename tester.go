package tester

import (
	"reflect"
	"strconv"
)

// DAO for database access object to define common behavior
type DAO interface {
	ExecuteStatements(setupStatements, teardownStatements, statements []Statement) ([]Result, []error, error)
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
	Actual   Result
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
	dao DAO,
	statements, setupStatements, teardownStatements, solutions []Statement,
	selectedQuestions []string,
) ([]TestResult, error) {
	var testResult = []TestResult{}

	i := 0
	results, errs, err := dao.ExecuteStatements(setupStatements, teardownStatements, statements)
	solutionResults, _, err := dao.ExecuteStatements(setupStatements, teardownStatements, solutions)
	testcases := ConvertTablesToTestCases(solutionResults)

	if err != nil {
		return testResult, err
	}

	for _, expected := range testcases {
		if !stringInSlice(expected.Index, selectedQuestions) && len(selectedQuestions) > 0 {
			continue
		}
		if i >= len(results) {
			testResult = append(
				testResult,
				TestResult{
					Expected: expected,
					Actual:   Result{},
					Pass:     false,
				},
			)
			i++
			continue
		}
		table := results[i]
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

// ConvertTablesToTestCases is a util method to add index to table
func ConvertTablesToTestCases(tables []Result) []TestCase {
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
