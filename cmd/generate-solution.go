package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/rcliao/sql-unit-test/db"
	"github.com/rcliao/sql-unit-test/parser"
)

func main() {
	sqlDB := getDBForSolution()

	solutionContent, err := ioutil.ReadFile("./solution.sql")
	if err != nil {
		panic(err)
	}
	statements := parser.ParseSQL(string(solutionContent), "#")
	var solution = make(map[string][]map[string]string)
	for i, statement := range statements {
		result, err := db.Query(sqlDB, statement.Text)
		if err != nil {
			panic(err)
		}
		solution[strconv.Itoa(i+1)] = result.Content
	}
	solutionJSON, _ := json.Marshal(solution)
	err = ioutil.WriteFile("./testcase.json", solutionJSON, 0644)
	fmt.Printf("%+v", solution)
}

func getDBForSolution() *sql.DB {
	// TODO: generate a random DB string
	defaultProtocol := "tcp"
	defaultPort := "3306"
	sqlDSN := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		defaultProtocol,
		os.Getenv("MYSQL_HOST"),
		defaultPort,
		os.Getenv("MYSQL_DB"),
	)

	db, err := sql.Open("mysql", sqlDSN)
	if err != nil {
		panic(err)
	}

	return db
}
