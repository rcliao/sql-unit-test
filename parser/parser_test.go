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
