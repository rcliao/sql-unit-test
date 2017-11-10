package db

import (
	"database/sql"

	mgo "gopkg.in/mgo.v2"

	tester "github.com/rcliao/sql-unit-test"
)

// Factory pattern for creating DAO
type Factory struct {
	sqlDB      *sql.DB
	mgoSession *mgo.Session
}

// NewFactory is constructor
func NewFactory(sqlDB *sql.DB, mgoSession *mgo.Session) *Factory {
	return &Factory{sqlDB, mgoSession}
}

// CreateDAO retruns the type of dao based on daoType (sql or mongo)
func (f *Factory) CreateDAO(daoType string) tester.DAO {
	if daoType == "mongo" {
		return NewMongoDAO(f.mgoSession)
	}
	return NewSQLDAO(f.sqlDB)
}
