package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"

	tester "github.com/rcliao/sql-unit-test"
	"github.com/rcliao/sql-unit-test/parser"
	"github.com/rcliao/sql-unit-test/runner"
)

var (
	configFilePath     = flag.String("config", "./config.json", "config.json filepath")
	testCaseFilePath   = flag.String("textcase", "./testcase.json", "textcase.json filepath")
	submissionFilePath = flag.String("submission", "./submission.txt", "submission.txt filepath")
)

func main() {
	flag.Parse()

	configContent, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		panic(err)
	}
	_, err = ioutil.ReadFile(*testCaseFilePath)
	if err != nil {
		panic(err)
	}
	_, err = ioutil.ReadFile(*submissionFilePath)
	if err != nil {
		panic(err)
	}

	config, err := parser.ParseConfig(string(configContent))
	if err != nil {
		panic(err)
	}
	db := getDB(config)
	runner := runner.NewMySQLRunner(db)

	result, _ := runner.Query("SELECT * FROM artists;")
	fmt.Println(result)
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
