package tester

// Submission represents each answer set student submit
type Submission struct {
	Question string
	Index    int
	Command  string
}

// Config for database connections
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
}
