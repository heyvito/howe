package disks

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/fatih/color"

	"github.com/victorgama/howe/helpers"
	"github.com/victorgama/howe/widgets"
)

func handle(payload map[string]interface{}, output chan interface{}, wait *sync.WaitGroup) {
	rawDisks, ok := payload["disks"]
	if !ok {
		output <- fmt.Errorf("disks: disks list not declared")
		wait.Done()
		return
	}

	disks := []string{}
	if diskArr, ok := rawDisks.([]interface{}); ok {
		for i, c := range diskArr {
			if n, ok := c.(string); ok {
				disks = append(disks, n)
			} else {
				output <- fmt.Errorf("disks: item %d in disks should be a string", i)
				wait.Done()
				return
			}
		}
	} else {
		output <- fmt.Errorf("disks: disks must be a list of strings")
		wait.Done()
		return
	}

	fsList := sigar.FileSystemList{}
	err := fsList.Get()
	if err != nil {
		helpers.ReportError(fmt.Sprintf("disks: %s", err))
		output <- "Could not read disk information"
		wait.Done()
		return
	}

	processableDisks := []fsItem{}
	if len(disks) == 1 && disks[0] == "*" {
		for _, fs := range fsList.List {
			result := fsItem{
				devName: fs.DevName,
				found:   true,
			}
			usage := sigar.FileSystemUsage{}
			if err := usage.Get(fs.DirName); err != nil {
				helpers.ReportError(fmt.Sprintf("disks: (%s) %s", fs.DirName, err))
				output <- "Error reading disks information"
				wait.Done()
				return
			}
			result.available = usage.Avail
			result.size = usage.Total
			result.used = usage.Used
			result.usePercent = usage.UsePercent()
			processableDisks = append(processableDisks, result)
		}
	} else {
		for _, fsName := range disks {
			result := fsItem{
				devName: fsName,
				found:   false,
			}
			for _, fs := range fsList.List {
				if fs.DevName == fsName {
					result.found = true
					usage := sigar.FileSystemUsage{}
					if err := usage.Get(fs.DirName); err != nil {
						helpers.ReportError(fmt.Sprintf("disks: %s", err))
						output <- "Error reading disks information"
						wait.Done()
						return
					}
					result.available = usage.Avail
					result.size = usage.Total
					result.used = usage.Used
					result.usePercent = usage.UsePercent()
				}
			}
			processableDisks = append(processableDisks, result)
		}
	}

	output <- process(processableDisks)
	wait.Done()
}

func formatSize(size uint64) string {
	return helpers.FormatSize(size * 1024)
}

func init() {
	widgets.Register("disks", handle)
}

type fsItem struct {
	found      bool
	devName    string
	size       uint64
	used       uint64
	available  uint64
	usePercent float64
}

func process(items []fsItem) string {
	nameColumnSize := 0
	for _, i := range items {
		l := len(i.devName)
		if l > nameColumnSize {
			nameColumnSize = l
		}
	}
	nameColumnSizeWithPadding := int(math.Min(13, float64(nameColumnSize+6)))
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	fmt.Fprint(w, padRight("Filesystems", " ", nameColumnSizeWithPadding))
	fmt.Fprint(w, "  Size  Used  Free  Use%\n")
	for _, fs := range items {
		if !fs.found {
			fmt.Fprintf(w, "  %s not found\n", fs.devName)
			continue
		}
		percentUse := int(math.Round(fs.usePercent))
		fmt.Fprint(w, "  "+padRight(fs.devName, " ", nameColumnSizeWithPadding))
		fmt.Fprint(w, padLeft(formatSize(fs.size), " ", 4)+"  ")
		fmt.Fprint(w, padLeft(formatSize(fs.used), " ", 4)+"  ")
		fmt.Fprint(w, padLeft(formatSize(fs.available), " ", 4)+"  ")
		fmt.Fprint(w, percentColor(percentUse).SprintFunc()(padLeft(strconv.Itoa(percentUse), " ", 3)+"%"))
		fmt.Fprint(w, "\n")

		totalSize := nameColumnSizeWithPadding + 20
		usedSize := int(math.Round(float64(totalSize) * (float64(percentUse) / float64(100))))

		fmt.Fprint(w, "  [")
		fmt.Fprint(w, percentColor(percentUse).SprintFunc()(strings.Repeat("=", usedSize)))
		fmt.Fprint(w, color.New(color.FgHiBlack).SprintFunc()(strings.Repeat("=", totalSize-usedSize)))
		fmt.Fprint(w, "]\n")
	}

	w.Flush()
	return buf.String()
}

func padLeft(str, fill string, size int) string {
	l := len(str)
	if l >= size {
		return str
	}
	return strings.Repeat(fill, size-l) + str
}

func padRight(str, fill string, size int) string {
	l := len(str)
	if l >= size {
		return str
	}
	return str + strings.Repeat(fill, size-l)
}

func percentColor(use int) *color.Color {
	c := color.FgWhite
	if use >= 0 && use < 75 {
		c = color.FgGreen
	} else if use >= 75 && use < 95 {
		c = color.FgYellow
	} else if use >= 95 {
		c = color.FgRed
	}
	return color.New(c)
}
