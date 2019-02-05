// +build !linux

package disks

// getMounts should return a map of mount point->device, but actually we're
// only supporting Linux here, so just return an empty map until we think
// what to do with this on other platforms.
func getMounts() (mountpoints, error) {
	return mountpoints{}, nil
}
