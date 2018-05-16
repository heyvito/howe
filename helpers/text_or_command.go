package helpers

import (
	"fmt"
	"os/exec"
)

// TextOrCommand interprets an input payload in order to determine whether it
// contains a static string, an command based on a string, or a command based on
// an array of arguments
func TextOrCommand(module string, payload map[string]interface{}) (string, error) {
	var toWrite string
	txt, ok := payload["text"]
	if !ok {
		cmd, ok := payload["command"]
		if !ok {
			return "", fmt.Errorf("%s: please provide text or command parameter", module)
		}

		var command *exec.Cmd

		switch cmdOrArr := cmd.(type) {
		case string:
			command = exec.Command(cmdOrArr)
		case []interface{}:
			if len(cmdOrArr) < 1 {
				return "", fmt.Errorf("%s: command array is empty", module)
			}
			cmds := []string{}
			for i, s := range cmdOrArr {
				if str, ok := s.(string); ok {
					cmds = append(cmds, str)
				} else {
					return "", fmt.Errorf("%s: item in position %d in command array should be a string", module, i)
				}
			}
			command = exec.Command(cmds[0], cmds[1:]...)
		}
		res, err := command.Output()
		if err != nil {
			return "", fmt.Errorf("%s: command failed: %s", module, err)
		}
		toWrite = string(res)
	} else {
		toWrite = txt.(string)
	}
	return toWrite, nil
}
