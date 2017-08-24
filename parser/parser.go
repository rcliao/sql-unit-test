package parser

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	tester "github.com/rcliao/sql-unit-test"
)

/*
ParseSQL parses through the content of text file to a list of submissions

A sample SQL text file may look like below

```
# 1. Find out all artist
SELECT *
FROM artists;

# 2. Find out all the song names
SELECT name
FROM songs;
```
*/
func ParseSQL(content, commentChar string) []tester.Statement {
	index := 1
	result := []tester.Statement{}
	submission := tester.Statement{}
	submission.Index = index
	commands := []string{}
	questions := []string{}
	lines := strings.Split(content, "\n")

	for i, l := range lines {
		// trim out of the unnecessary tailing and opening spaces
		line := strings.Trim(l, " ")
		numberOfCharacters := utf8.RuneCountInString(line)
		// when face empty line, add submission and reset state back to start
		if numberOfCharacters == 0 {
			if len(commands) == 0 && len(questions) == 0 {
				// means the previous line is also empty, if so, don't do anything
				continue
			}
			// add the submission to result list
			index++
			submission.Index = index
			submission.Text = strings.Join(commands, " ")
			submission.Comment = strings.Join(questions, " ")
			result = append(result, submission)

			// reset all the state
			submission = tester.Statement{}
			commands = []string{}
			questions = []string{}
			continue
		}

		if strings.HasPrefix(line, commentChar) {
			questions = append(questions, line)
		} else {
			commands = append(commands, line)
		}
		if i == len(lines)-1 || line[len(line)-1:] == ";" {
			// add the submission to result list
			if len(commands) > 0 {
				index++
				submission.Index = index
				submission.Text = strings.Join(commands, " ")
				submission.Comment = strings.Join(questions, " ")
				result = append(result, submission)
				// reset
				submission = tester.Statement{}
				commands = []string{}
				questions = []string{}
			}
		}
	}
	return result
}

/*
ParseTestCases parses through content defined in JSON for test cases

An usual test cases may be defined as below:

```
{
	1: [
		{
			"name": "Eric",
			"age": "21"
		}
	]
}
```
*/
func ParseTestCases(content string) (map[string][]map[string]string, error) {
	result := map[string][]map[string]string{}
	decoder := json.NewDecoder(strings.NewReader(content))
	if err := decoder.Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

// ParseConfig parses the config file in the json format to type object
func ParseConfig(content string) (tester.Config, error) {
	result := tester.Config{}
	decoder := json.NewDecoder(strings.NewReader(content))
	if err := decoder.Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}
