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
)

func main() {
	flag.Parse()

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

	config, err := parser.ParseConfig(string(configContent))
	if err != nil {
		panic(err)
	}
	db := getDB(config)
	runner := runner.NewMySQLRunner(db)

	submissions := parser.ParseSQLSubmission(string(submissionContent), "#")
	testCases, err := parser.ParseTestCases(string(testCasesContent))
	if err != nil {
		panic(err)
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
