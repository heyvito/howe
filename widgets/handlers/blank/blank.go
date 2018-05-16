package blank

import (
	"sync"

	"github.com/victorgama/howe/widgets"
)

func handle(payload map[string]interface{}, output chan interface{}, wait *sync.WaitGroup) {
	output <- " "
	wait.Done()
}

func init() {
	widgets.Register("blank", handle)
}
