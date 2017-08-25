package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	tester "github.com/rcliao/sql-unit-test"
	"github.com/rcliao/sql-unit-test/parser"
	"github.com/rcliao/sql-unit-test/runner"
)

var (
	configFilePath     = flag.String("config", "./config.json", "config.json filepath")
	testCaseFilePath   = flag.String("testcases", "./testcase.json", "testcase.json filepath")
	submissionFilePath = flag.String("submission", "./statements.sql", "submission.txt filepath")
	// optional params
	setupFilePath    = flag.String("setup", "", "setup.sql filepath")
	teardownFilePath = flag.String("teardown", "", "teardown.sql filepath")
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

	db := getDB(config)
	runner := runner.NewMySQLRunner(db)
	failedTestCases, err := tester.Run(runner, submissions, setupStatements, teardownStatements, testCases)
	if err != nil {
		panic(err)
	}
	if len(failedTestCases) == 0 {
		fmt.Println("All test passed!")
	} else {
		fmt.Println(strings.Join(failedTestCases, "\n"))
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
