package parser

import "testing"

func TestParseSQLSubmissions(t *testing.T) {
	var content = `# 1. Find out all artist
SELECT * FROM artists;


# 2. Find out all the song names
SELECT name
FROM songs;

SELECT * FROM artists;
`
	submissions := ParseSQLSubmission(content, "#")

	if submissions[0].Command != "SELECT * FROM artists;" {
		t.Error("Expected 'first submission command' to be 'SELECT * FROM artists' got '", submissions[0].Command+"'")
	}

	if submissions[1].Command != "SELECT name FROM songs;" {
		t.Error("Expect 'second submission command' to be 'SELECT name FROM songs;' got '" + submissions[1].Command + "'")
	}

	if submissions[2].Command != "SELECT * FROM artists;" {
		t.Error("Expect 'third submission command' to be 'SELECT * FROM artists;' got '" + submissions[2].Command + "'")
	}
}

func TestParseSQLSubmissions2(t *testing.T) {
	var content = `# 1. Find out all artist
SELECT * FROM artists;`
	submissions := ParseSQLSubmission(content, "#")

	if submissions[0].Command != "SELECT * FROM artists;" {
		t.Error("Expected 'first submission command' to be 'SELECT * FROM artists' got '", submissions[0].Command+"'")
	}
}

func TestParseTestCases(t *testing.T) {
	var content = `{
	"1": [
		{
			"name": "Eric",
			"age": "21"
		}
	]
}`
	testCases, err := ParseTestCases(content)

	if err != nil {
		t.Error("Having error parsing the content", err)
	}

	if testCases["1"][0]["name"] != "Eric" {
		t.Error("Expect test case to get correct name 'Eric' but got '" + testCases["1"][0]["name"] + "'")
	}

	if testCases["1"][0]["age"] != "21" {
		t.Error("Expect test case to get correct age '21' but got '" + testCases["1"][0]["age"] + "'")
	}
}

func TestParseConfig(t *testing.T) {
	var content = `{
	"username": "root",
	"host": "localhost",
	"database": "cs1222"
}`
	result, err := ParseConfig(content)

	if err != nil {
		t.Error("Having error parseing the json content", err)
	}

	if result.Username != "root" {
		t.Error("Expect to get 'root' username but got '" + result.Username + "'")
	}

	if result.Password != "" {
		t.Error("Expect to get '' password but got '" + result.Password + "'")
	}

	if result.Host != "localhost" {
		t.Error("Expect to get 'localhost' host but got '" + result.Host + "'")
	}

	if result.Database != "cs1222" {
		t.Error("Expect to get 'cs1222' host but got '" + result.Database + "'")
	}
}
