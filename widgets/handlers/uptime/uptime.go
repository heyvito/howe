package uptime

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cloudfoundry/gosigar"

	"github.com/victorgama/howe/helpers"
	"github.com/victorgama/howe/widgets"
)

func handle(payload map[string]interface{}, output chan interface{}, wait *sync.WaitGroup) {
	uptime := sigar.Uptime{}
	err := uptime.Get()
	if err != nil {
		helpers.ReportError(fmt.Sprintf("uptime: %s", err))
		output <- "uptime: No information available"
		wait.Done()
		return
	}

	time := uint64(uptime.Length)

	components := []string{}

	days := time / (60 * 60 * 24)

	if days != 0 {
		s := ""
		if days > 1 {
			s = "s"
		}
		components = append(components, fmt.Sprintf("%d day%s", days, s))
	}

	minutes := time / 60
	hours := minutes / 60
	hours %= 24
	minutes %= 60

	if hours > 0 {
		s := ""
		if hours > 1 {
			s = "s"
		}

		components = append(components, fmt.Sprintf("%d hour%s", hours, s))
	}

	if minutes > 0 {
		s := ""
		if minutes > 1 {
			s = "s"
		}
		components = append(components, fmt.Sprintf("%d minute%s", minutes, s))
	}

	output <- fmt.Sprintf("up %s", strings.Join(components, ", "))
	wait.Done()
}

func init() {
	widgets.Register("uptime", handle)
}
