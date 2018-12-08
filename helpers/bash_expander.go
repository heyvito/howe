package helpers

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Based on IEEE Std 1003.1-2008 / IEEE POSIX P1003.2/ISO 9945.2
// See also: http://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_10_02
// The next expression matches $TEST, and \$TEST. Escaping must be
var environmentVariableMatcher = regexp.MustCompile(`(?im)(?:(\\)?\$)([a-zA-Z_]+[a-zA-Z0-9_]*)`)

var runCommandMatcher = regexp.MustCompile(`(?im)(\\)?\$\(((?:[^\\(\\)]+|\.)*)\)`)

// ExpandBashStyleString expands variables ($TERM), replacing unknown ones with
// an empty string, and running any commands wrapped in "$()". If running a
// command fails, and error is returned instead.
func ExpandBashStyleString(input string) (string, error) {
	result, err := replaceAllStringSubmatchFunc(environmentVariableMatcher, input, expandEnvironmentVariables)
	if err != nil {
		return "", err
	}

	return replaceAllStringSubmatchFunc(runCommandMatcher, result, expandRunCommandExpression)
}

func expandEnvironmentVariables(match []string) (string, error) {
	if len(match) != 3 {
		return "", nil
	}
	if match[1] == "\\" {
		return match[0][1:], nil
	}

	return os.Getenv(match[2]), nil
}

func expandRunCommandExpression(match []string) (string, error) {
	if len(match) != 3 {
		return "", nil
	}

	if match[1] == "\\" {
		return match[0][1:], nil
	}

	args := strings.Split(match[2], " ")
	command := exec.Command(args[0], args[1:]...)
	result, err := command.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(result)), nil
}

func replaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) (string, error)) (string, error) {
	groups := re.FindAllStringSubmatchIndex(str, -1)
	if groups == nil {
		return str, nil
	}

	delta := 0
	result := str
	for _, match := range groups {
		matches := []string{}
		for i := 0; i < len(match); i += 2 {
			if match[i] == -1 && match[i+1] == -1 {
				matches = append(matches, "")
			} else {
				matches = append(matches, str[match[i]:match[i+1]])
			}
		}

		lastStrLength := match[1] - match[0]
		replacement, err := repl(matches)
		if err != nil {
			return "", err
		}
		result = result[:match[0]+delta] + replacement + result[match[1]+delta:]
		delta += len(replacement) - lastStrLength
	}

	return result, nil
}
