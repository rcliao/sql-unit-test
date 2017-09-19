# SQL Unit Test

A GoLang Web application to run SQL unit tests defined in JSON.

## Get Started

This repository uses [dep](https://github.com/golang/dep) as dependency management
tool. That is, you can use `dep ensure` to install all the dependencies (MySQL
Driver for example).

After installing the runtime dependencies, you will need to ensure the MySQL
server is running so that you can connect to it for running tests.

To configure application to connect to MySQL, follow the following environment
variables:

```
export MYSQL_USERNAME=root
export MYSQL_PASSWORD=
export MYSQL_HOST=localhost
```

Then, you can run application by `go run cmd/server.go` to see application
running at http://localhost:8000

## Ideas

- [ ] Add control panel for defining new test cases
- [x] Define setup and teardown life cycle
- [x] Add web interface for ease of distribution
- [x] Use UUID (or some other random string) for database name for parallel testings
- [x] Improve the existing front end to have better style/experience
    - [x] Add Code Editor
    - [x] Better feedback on fail test cases
- [x] Use dep to do dependency management

## Roadmap

- [ ] Add other SQL driver implementations
- [ ] Add prepared statement support
