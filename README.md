# SQL Unit Test

A GoLang CLI app to run SQL unit tests.

## Get Started

Before get started, you will need to install "github.com/go-sql-driver/mysql"
by `go get github.com/go-sql-driver/mysql` Assuming you are writing unit tests
for MySQL.

After installing the specific database driver, the other dependencies left is
`config.json`, `testcase.json` and `submission.sql`.

### config.json

Config.json defines the necessary configuration to connect to database.

A sample of config.json can be found below:

```json
{
    "username": "root",
    "host": "localhost",
    "database": "lyrics"
}
```

### submission.sql

Submission defines the sql submission.

```
# 1. Find the alternative genre
SELECT Genre FROM genre where Genre = 'alternative';
```

### testcase.json

Test cases defines the index to expected outcome.

Test

```json
{
    "1": [
        {
            "genre": "alternative"
        }
    ]
}
```

### Running tests

To run tests, `go run cmd/main.go`

> You may provide the path to each of the files above through command line flags (testcase, submission & config)

## Ideas

- [x] Define setup and teardown life cycle
- [ ] Add web interface for ease of distribution
- [ ] Add other SQL driver implementations
