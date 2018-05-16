package helpers

import (
	"bufio"
	"bytes"
	"fmt"
)

// FormatSize formats a given size into K/M/G/T/P/E
func FormatSize(size uint64) string {
	ord := []string{"K", "M", "G", "T", "P", "E"}
	o := 0
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)

	if size < 973 {
		fmt.Fprintf(w, "%3d ", size)
		w.Flush()
		return buf.String()
	}

	for {
		remain := size & 1023
		size >>= 10

		if size >= 973 {
			o++
			continue
		}

		if size < 9 || (size == 9 && remain < 973) {
			remain = ((remain * 5) + 256) / 512
			if remain >= 10 {
				size++
				remain = 0
			}

			fmt.Fprintf(w, "%d.%d%s", size, remain, ord[o])
			break
		}

		if remain >= 512 {
			size++
		}

		fmt.Fprintf(w, "%3d%s", size, ord[o])
		break
	}

	w.Flush()
	return buf.String()
}
