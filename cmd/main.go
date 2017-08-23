package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	tester "github.com/rcliao/sql-unit-test"
	"github.com/rcliao/sql-unit-test/parser"
	"github.com/rcliao/sql-unit-test/runner"
)

var (
	configFilePath     = flag.String("config", "./config.json", "config.json filepath")
	testCaseFilePath   = flag.String("testcases", "./testcase.json", "testcase.json filepath")
	submissionFilePath = flag.String("submission", "./submission.sql", "submission.txt filepath")
	setupFilePath      = flag.String("setup", "", "setup.sql filepath")
	teardownFilePath   = flag.String("teardown", "", "teardown.sql filepath")
)

func main() {
	flag.Parse()

	// getting contents
	configContent, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		panic(err)
	}
	testCasesContent, err := ioutil.ReadFile(*testCaseFilePath)
	if err != nil {
		panic(err)
	}
	submissionContent, err := ioutil.ReadFile(*submissionFilePath)
	if err != nil {
		panic(err)
	}
	var setupContent = []byte{}
	if *setupFilePath != "" {
		setupContent, err = ioutil.ReadFile(*setupFilePath)
		if err != nil {
			panic(err)
		}
	}
	var teardownContent = []byte{}
	if *teardownFilePath != "" {
		teardownContent, err = ioutil.ReadFile(*teardownFilePath)
		if err != nil {
			panic(err)
		}
	}

	// parsing content
	config, err := parser.ParseConfig(string(configContent))
	if err != nil {
		panic(err)
	}
	submissions := parser.ParseSQL(string(submissionContent), "#")
	testCases, err := parser.ParseTestCases(string(testCasesContent))
	if err != nil {
		panic(err)
	}
	var setupStatements = []tester.Submission{}
	if string(setupContent) != "" {
		setupStatements = parser.ParseSQL(string(setupContent), "#")
		if err != nil {
			panic(err)
		}
	}
	var teardownStatements = []tester.Submission{}
	if string(teardownContent) != "" {
		teardownStatements = parser.ParseSQL(string(teardownContent), "#")
		if err != nil {
			panic(err)
		}
	}

	db := getDB(config)
	runner := runner.NewMySQLRunner(db)
	// ready to run through life cycle
	// setup
	if len(setupStatements) > 0 {
		for _, statement := range setupStatements {
			fmt.Println(statement.Command)
			if err := runner.Execute(statement.Command); err != nil {
				panic(err)
			}
		}
	}
	var pass = true
	for i, submission := range submissions {
		result, err := runner.Query(submission.Command)
		if err != nil {
			panic(err)
		}
		expected := testCases[strconv.Itoa(i+1)]
		if !reflect.DeepEqual(result, expected) {
			fmt.Printf("Test case %d: expected %v but got %v", i, expected, result)
			pass = false
		}
	}
	if pass {
		fmt.Println("All tests passed!")
	}
	// teardown
	if len(teardownStatements) > 0 {
		for _, statement := range teardownStatements {
			if err := runner.Execute(statement.Command); err != nil {
				panic(err)
			}
		}
	}
}

func getDB(config tester.Config) *sql.DB {
	defaultProtocol := "tcp"
	defaultPort := "3306"
	sqlDSN := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s",
		config.Username,
		config.Password,
		defaultProtocol,
		config.Host,
		defaultPort,
		config.Database,
	)

	db, err := sql.Open("mysql", sqlDSN)
	if err != nil {
		panic(err)
	}

	return db
}
