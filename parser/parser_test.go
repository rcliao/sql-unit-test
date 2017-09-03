package parser

import (
	"fmt"
	"testing"
)

func TestParseSQL(t *testing.T) {
	var content = `# 1. Find out all artist
SELECT * FROM artists;


# 2. Find out all the song names
SELECT name
FROM songs;

SELECT * FROM artists;`
	submissions := ParseSQL(content, "#")

	if len(submissions) != 3 {
		t.Error(fmt.Sprint("Expected the submission size to be '3' but got", len(submissions)))
	}

	if submissions[0].Text != "SELECT * FROM artists;" {
		t.Error("Expected 'first submission command' to be 'SELECT * FROM artists' got '", submissions[0].Text+"'")
	}

	if submissions[1].Text != "SELECT name FROM songs;" {
		t.Error("Expect 'second submission command' to be 'SELECT name FROM songs;' got '" + submissions[1].Text + "'")
	}

	if submissions[2].Text != "SELECT * FROM artists;" {
		t.Error("Expect 'third submission command' to be 'SELECT * FROM artists;' got '" + submissions[2].Text + "'")
	}
}

func TestParseSQL2(t *testing.T) {
	var content = `# 1. Find out all artist
SELECT * FROM artists;`
	submissions := ParseSQL(content, "#")

	if submissions[0].Text != "SELECT * FROM artists;" {
		t.Error("Expected 'first submission command' to be 'SELECT * FROM artists' got '" + submissions[0].Text + "'")
	}
}

func TestParseSQL3(t *testing.T) {
	var content = `SELECT * FROM artists;
UPDATE artists SET Name = 'Eric' WHERE ID = 1;`
	submissions := ParseSQL(content, "#")

	if submissions[0].Text != "SELECT * FROM artists;" {
		t.Error("Expected 'first submission command' to be 'SELECT * FROM artists' got '" + submissions[0].Text + "'")
	}

	if submissions[1].Text != "UPDATE artists SET Name = 'Eric' WHERE ID = 1;" {
		t.Error("Expected 'second submission command' to be 'UPDATE artists SET Name = 'Eric' WHERE ID = 1;' got '", submissions[1].Text+"'")
	}
}

func TestParseSQL4(t *testing.T) {
	var content = `SELECT Title, UPC, Genre FROM Titles;

SELECT * FROM Titles WHERE ArtistID=2;

SELECT FirstName, LastName, HomePhone, EMail FROM Members;

SELECT MemberID FROM Members WHERE Gender = 'M';

SELECT MemberID, Country FROM Members WHERE Country = 'Canada';`
	submissions := ParseSQL(content, "#")
	if submissions[0].Text != "SELECT Title, UPC, Genre FROM Titles;" {
		t.Error("Expected first submission command to be 'SELECT Title, UPC, Genre FROM Titles;' but got '", submissions[0].Text)
	}
	if submissions[1].Text != "SELECT * FROM Titles WHERE ArtistID=2;" {
		t.Error("Expected first submission command to be 'SELECT * FROM Titles WHERE ArtistID=2;' but got '", submissions[1].Text)
	}
	if submissions[2].Text != "SELECT FirstName, LastName, HomePhone, EMail FROM Members;" {
		t.Error("Expected first submission command to be 'SELECT FirstName, LastName, HomePhone, EMail FROM Members;' but got '", submissions[2].Text)
	}
	if submissions[3].Text != "SELECT MemberID FROM Members WHERE Gender = 'M';" {
		t.Error("Expected first submission command to be 'SELECT MemberID FROM Members WHERE Gender = 'M';' but got '", submissions[3].Text)
	}
	if submissions[4].Text != "SELECT MemberID, Country FROM Members WHERE Country = 'Canada';" {
		t.Error("Expected first submission command to be 'SELECT MemberID, Country FROM Members WHERE Country = 'Canada'' but got '", submissions[4].Text)
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

	if testCases[0].Content[0]["name"] != "Eric" {
		t.Error("Expect test case to get correct name 'Eric' but got '" + testCases[0].Content[0]["name"] + "'")
	}

	if testCases[0].Content[0]["age"] != "21" {
		t.Error("Expect test case to get correct age '21' but got '" + testCases[0].Content[0]["age"] + "'")
	}
}

func TestParseConfig(t *testing.T) {
	var content = `{
	"username": "root",
	"host": "localhost"
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
}
