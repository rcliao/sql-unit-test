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

### statements.sql

Statements defines the sql submission to be tested in that order.

```
# 1. Find the alternative genre
SELECT Genre FROM genre where Genre = 'alternative';
```

> You can also provide the setup.sql and teardown.sql as same format as above to
> define better test environment control

### testcase.json

Test cases defines the index to expected outcome.

Test

```json
{
    "1": [
        {
            "Genre": "alternative"
        }
    ]
}
```

### Running tests

To run tests, `go run cmd/main.go`

> You may provide the path to each of the files above through command line flags (testcase, submission & config)

## Ideas

- [x] Define setup and teardown life cycle
- [x] Add web interface for ease of distribution
- [ ] Use UUID (or some other random string) for database name for parallel testings
- [ ] Improve the existing front end to have better style/experience
    - [ ] Add Code Editor
    - [ ] Better feedback on fail test cases
- [ ] Add control panel for defining new test cases
- [ ] Use dep to do dependency management

## Roadmap

- [ ] Add other SQL driver implementations
- [ ] Add prepared statement support
