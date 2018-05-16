package updates

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"

	"github.com/victorgama/howe/helpers"
	"github.com/victorgama/howe/widgets"
)

func handle(payload map[string]interface{}, output chan interface{}, wait *sync.WaitGroup) {
	result, err := exec.Command("sh", "-c", "apt-get -s -o Debug::NoLocking=true upgrade | grep ^Inst | wc -l").Output()

	if err != nil {
		helpers.ReportError(fmt.Sprintf("updates: %s", err))
		output <- "No update information available"
		wait.Done()
		return
	}

	number, err := strconv.Atoi(strings.TrimSpace(string(result)))
	if err != nil {
		helpers.ReportError(fmt.Sprintf("updates: %s", err))
		output <- "No update information available"
		wait.Done()
		return
	}

	if number == 0 {
		output <- "No updates available."
	} else {
		s := ""
		if number > 1 {
			s = "s"
		}
		output <- color.YellowString(fmt.Sprintf("%d update%s available", number, s))
	}
	wait.Done()
}

func init() {
	widgets.Register("updates", handle)
}
