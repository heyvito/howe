package docker

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/color"
	dockerApi "github.com/fsouza/go-dockerclient"

	"github.com/victorgama/howe/helpers"
	"github.com/victorgama/howe/widgets"
)

const defaultDockerSock = "unix:///var/run/docker.sock"

func handle(payload map[string]interface{}, output chan interface{}, wait *sync.WaitGroup) {
	rawContainers, ok := payload["containers"]
	if !ok {
		output <- fmt.Errorf("docker: containers list not declared")
		wait.Done()
		return
	}

	containers := []string{}
	if containerArr, ok := rawContainers.([]interface{}); ok {
		for i, c := range containerArr {
			if n, ok := c.(string); ok {
				containers = append(containers, n)
			} else {
				output <- fmt.Errorf("docker: item %d in containers should be a string", i)
				wait.Done()
				return
			}
		}
	} else {
		output <- fmt.Errorf("docker: containers must be a list of strings")
		wait.Done()
		return
	}

	dockerSock := defaultDockerSock

	if rawSocket, ok := payload["socket"]; ok {
		if socket, ok := rawSocket.(string); ok {
			dockerSock = socket
		}
	}

	docker, err := dockerApi.NewClient(dockerSock)
	if err != nil {
		helpers.ReportError(fmt.Sprintf("docker: %s", err))
		output <- fmt.Sprintf("Docker: Could not connect on %s", dockerSock)
		wait.Done()
		return
	}

	dockerContainers, err := docker.ListContainers(dockerApi.ListContainersOptions{All: true, Limit: 10000})
	if err != nil {
		helpers.ReportError(fmt.Sprintf("docker: %s", err))
		output <- fmt.Sprintf("Docker: Error enumerating containers.")
		wait.Done()
		return
	}

	results := [][]string{}

	for _, n := range containers {
		if strings.HasPrefix(n, "regexp:") {
			findContainersRegex(&dockerContainers, n[7:], &results)
		} else {
			findContainersStatic(&dockerContainers, n, &results)
		}
	}

	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	longest := longestString(results)

	for _, v := range results {
		fmt.Fprintf(w, "    %s    %s\n", padString(v[0], longest), v[1])
	}
	w.Flush()

	output <- "\nDocker:\n" + buf.String()
	wait.Done()
}

func init() {
	widgets.Register("docker", handle)
}

func longestString(list [][]string) int {
	longest := 0
	for _, s := range list {
		l := len(s[0])
		if l > longest {
			longest = l
		}
	}

	return longest
}

func findContainersStatic(containers *[]dockerApi.APIContainers, name string, results *[][]string) {
	r := "not found"
	f := color.FgRed
	for _, c := range *containers {
		found := false
		for _, cn := range c.Names {
			if strings.Replace(cn, "/", "", -1) == name {
				found = true
				break
			}
		}

		if found {
			r = helpers.Titleize(c.State) + ", " + c.Status
			switch strings.ToLower(c.State) {
			case "paused", "created":
				f = color.FgWhite
			case "restarting":
				f = color.FgYellow
			case "running":
				f = color.FgGreen
			case "dead", "exited":
				f = color.FgRed
			}
		}
	}
	tmpResults := append(*results, []string{name, color.New(f).SprintFunc()(r)})
	*results = tmpResults
}

func findContainersRegex(containers *[]dockerApi.APIContainers, name string, results *[][]string) {
	r := "not found"
	f := color.FgRed
	re, err := regexp.Compile(name)
	if err != nil {
		tmpResults := append(*results, []string{name, fmt.Sprintf("Error compiling regexp: %s", err)})
		*results = tmpResults
		return
	}

	for _, c := range *containers {
		found := false
		cname := ""
		for _, cn := range c.Names {
			currentName := strings.Replace(cn, "/", "", -1)
			if re.Match([]byte(currentName)) {
				found = true
				cname = currentName
				break
			}
		}
		if !found {
			continue
		}

		if found {
			r = helpers.Titleize(c.State) + ", " + c.Status
			switch strings.ToLower(c.State) {
			case "paused", "created":
				f = color.FgWhite
			case "restarting":
				f = color.FgYellow
			case "running":
				f = color.FgGreen
			case "dead", "exited":
				f = color.FgRed
			}
		}
		tmpResults := append(*results, []string{cname, color.New(f).SprintFunc()(r)})
		*results = tmpResults
	}

}

func padString(str string, size int) string {
	strLen := len(str)
	if strLen >= size {
		return str + ":"
	}
	return str + ":" + strings.Repeat(" ", size-strLen)
}
