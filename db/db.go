package db

import "database/sql"

// Table contains the MySQL Table result in map of string to string
type Table struct {
	// IDEA: probably think of a better way to do type comparison
	Query   string
	Content []map[string]string
}

// Query a query and return the table
func Query(db *sql.Tx, query string) (Table, error) {
	result := Table{Query: query}

	rows, err := db.Query(query)
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
