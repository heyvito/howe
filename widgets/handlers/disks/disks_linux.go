// +build linux

package disks

import (
	"bufio"
	"io"
	"os"
	"regexp"
)

var mountPointRegexp = regexp.MustCompile(`^(\/[^\s]*)[\s\t]+(\/[^\s]*)`)

// getMounts returns a map of mount point->device
func getMounts() (mountpoints, error) {
	mounts := mountpoints{}
	f, err := os.OpenFile("/proc/mounts", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
		opts := mountPointRegexp.FindAllStringSubmatch(line, -1)
		if opts != nil {
			for _, match := range opts {
				var (
					device = match[1]
					mount  = match[2]
				)
				mounts = append(mounts, mountpoint{device, mount})
			}
		}
	}
	return mounts, nil
}
