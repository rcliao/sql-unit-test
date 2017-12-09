package db

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/pkg/errors"
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

	for _, statement := range setupStatements {
		commands := strings.Split(strings.Replace(statement.Text, "{dbName}", randomDatabaseName, -1), " ")
		cmd := exec.Command(commands[0], commands[1:]...)
		err := cmd.Run()
		if err != nil {
			fmt.Println("has error running setup statement", err)
		}
	}

	var cErr error
	for _, statement := range statements {
		r := bson.M{}
		statementResult := tester.Result{Query: statement.Text}
		err := m.session.DB(randomDatabaseName).Run(bson.M{"eval": statement.Text}, &r)
		if err != nil {
			if cErr != nil {
				cErr = errors.Wrap(cErr, err.Error())
			} else {
				cErr = err
			}
		}
		if strings.Contains(statement.Text, "count()") {
			errs = append(errs, cErr)
			contentRaw, okay := r["retval"]
			if !okay || cErr != nil {
				m := map[string]string{}
				statementResult.Content = append(statementResult.Content, m)
				result = append(result, statementResult)
				continue
			}
			content := contentRaw.(float64)
			m := map[string]string{}
			m["0"] = strconv.FormatFloat(content, 'f', -1, 64)
			statementResult.Content = append(statementResult.Content, m)
			result = append(result, statementResult)
		}
		if strings.Contains(statement.Text, "aggregate") {
			errs = append(errs, cErr)
			contentRaw, okay := r["retval"]
			if !okay || cErr != nil {
				m := map[string]string{}
				statementResult.Content = append(statementResult.Content, m)
				result = append(result, statementResult)
				continue
			}
			content := contentRaw.(bson.M)["_batch"].([]interface{})
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
			statementResult.Content = reverse(statementResult.Content)
			result = append(result, statementResult)
		}
		if strings.Contains(statement.Text, "toArray()") {
			errs = append(errs, cErr)
			contentRaw, okay := r["retval"]
			if !okay || cErr != nil {
				m := map[string]string{}
				statementResult.Content = append(statementResult.Content, m)
				result = append(result, statementResult)
				continue
			}
			content := contentRaw.([]interface{})
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

func reverse(list []map[string]string) []map[string]string {
	s := make([]map[string]string, len(list))
	for i := 0; i < len(list); i++ {
		s[i] = list[i]
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
