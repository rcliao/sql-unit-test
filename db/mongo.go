package db

import (
	"fmt"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	tester "github.com/rcliao/sql-unit-test"
)

// MongoDAO implements DAO interface to interact with Mongo
type MongoDAO struct {
	session *mgo.Session
}

// NewMongoDAO is constructor
func NewMongoDAO(session *mgo.Session) *MongoDAO {
	return &MongoDAO{session}
}

// ExecuteStatements builds the statements into result (checking if result are correct)
func (m *MongoDAO) ExecuteStatements(setupStatements, teardownStatements, statements []tester.Statement) ([]tester.Result, []error, error) {
	result := []tester.Result{}
	errs := []error{}

	randomDatabaseName := getRandomString()

	for _, statement := range statements {
		r := bson.M{}
		statementResult := tester.Result{Query: statement.Text}
		err := m.session.DB(randomDatabaseName).Run(bson.M{"eval": statement.Text}, &r)
		if strings.Contains(statement.Text, "toArray()") {
			errs = append(errs, err)
			content := r["retval"].([]interface{})
			for _, record := range content {
				recordContent := record.(bson.M)
				m := map[string]string{}
				for k, v := range recordContent {
					if k != "_id" {
						m[k] = fmt.Sprintf("%v", v)
					}
				}
				statementResult.Content = append(statementResult.Content, m)
			}
			result = append(result, statementResult)
		}
	}
	// cleanup
	m.session.DB(randomDatabaseName).Run(bson.M{"eval": "db.dropDatabase();"}, nil)

	return result, errs, nil
}
