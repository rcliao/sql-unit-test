package db

import (
	"database/sql"
	"log"
	"strings"

	"github.com/pkg/errors"
	tester "github.com/rcliao/sql-unit-test"
)

// SQLDAO implements DAO interface for SQL
type SQLDAO struct {
	sqlDB *sql.DB
}

// NewSQLDAO is constructor pattern
func NewSQLDAO(sqlDB *sql.DB) *SQLDAO {
	return &SQLDAO{sqlDB}
}

// ExecuteStatements takes arguments to build a list of tables
func (d *SQLDAO) ExecuteStatements(setupStatements, teardownStatements, statements []tester.Statement) ([]tester.Result, []error, error) {
	result := []tester.Result{}
	errs := []error{}

	// use random 32 characters database
	randomDatabaseName := getRandomString()

	tx, err := d.sqlDB.Begin()
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
			table, err := Query(tx, statement.Text)
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

// Query a query and return the table
func Query(tx *sql.Tx, query string) (tester.Result, error) {
	result := tester.Result{Query: query}

	rows, err := tx.Query(query)
	if err != nil {
		return result, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return result, err
	}
	// Result is your slice string.
	rawResult := make([][]byte, len(columns))
	dest := make([]interface{}, len(columns)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}
	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			return result, err
		}
		tableRow := map[string]string{}
		for i, raw := range rawResult {
			val := ""
			if raw != nil {
				val = string(raw)
			}
			tableRow[columns[i]] = val
		}
		result.Content = append(result.Content, tableRow)
	}

	return result, nil
}
