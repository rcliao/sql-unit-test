package runner

import "database/sql"

// Runner defines behavior to execute query and return table rows
type Runner interface {
	Query(query string) ([]map[string]string, error)
	Execute(query string) error
}

// MySQLRunner implements Runner for MySQL specification
type MySQLRunner struct {
	db *sql.DB
}

// NewMySQLRunner is a simple constructor pattern for getting MySQLRunner
func NewMySQLRunner(db *sql.DB) MySQLRunner {
	return MySQLRunner{db}
}

// Query a query and return the table
func (r MySQLRunner) Query(query string) ([]map[string]string, error) {
	result := []map[string]string{}

	rows, err := r.db.Query(query)
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
		for _, raw := range rawResult {
			val := ""
			if raw != nil {
				val = string(raw)
			}
			tableRow := map[string]string{}
			for _, col := range columns {
				tableRow[col] = val
			}
			result = append(result, tableRow)
		}
	}

	return result, nil
}

// Execute a query to update database
func (r MySQLRunner) Execute(query string) error {
	_, err := r.db.Exec(query)
	return err
}
