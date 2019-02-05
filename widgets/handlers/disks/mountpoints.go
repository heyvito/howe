package disks

type mountpoint struct {
	device     string
	mountpoint string
}

type mountpoints []mountpoint

func (m mountpoints) byDevice(name string) *mountpoint {
	for _, i := range m {
		if i.device == name {
			return &i
		}
	}
	return nil
}

func (m mountpoints) byMountpoint(path string) *mountpoint {
	for _, i := range m {
		if i.mountpoint == path {
			return &i
		}
	}
	return nil
}
