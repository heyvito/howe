package print

import (
	"sync"

	"github.com/victorgama/howe/helpers"
	"github.com/victorgama/howe/widgets"
)

func handle(payload map[string]interface{}, output chan interface{}, wait *sync.WaitGroup) {
	toWrite, err := helpers.TextOrCommand("print", payload)
	if err != nil {
		output <- err
		wait.Done()
		return
	}
	output <- toWrite
	wait.Done()
}

func init() {
	widgets.Register("print", handle)
}
